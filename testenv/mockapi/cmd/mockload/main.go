// Copyright 2026 Haikei Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command mockload drives the contract mock at a fixed request rate and
// reports latency percentiles and the error rate. It exercises both
// transports with generated operations: GitHub over REST, Linear over
// GraphQL.
//
//	mockload -addr :8080 -rps 50 -duration 10s
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	githubv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1/githubops"
	linearv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1/linearops"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

func main() {
	addr := flag.String("addr", "http://127.0.0.1:8080", "mockapi base URL")
	rps := flag.Int("rps", 50, "requests per second")
	duration := flag.Duration("duration", 10*time.Second, "how long to drive load")
	timeout := flag.Duration("timeout", 5*time.Second, "per-request timeout")
	seed := flag.Int("seed", 1, "seed for the synthetic ids sent")
	flag.Parse()
	conn, err := connection.NewHTTP(*addr, &contracts.Connection{Name: "pat", Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_BEARER}, "mock-load-token", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}
	github := githubops.New(conn)
	linear := linearops.New(conn)

	deadline := time.Now().Add(*duration)
	tick := time.Second / time.Duration(*rps)
	perRequest := 2 // github + linear call per tick
	total := *rps * int(duration.Seconds()) * perRequest
	latencies := make([]time.Duration, 0, total)
	var mu sync.Mutex
	var failures int

	run := func(call func(ctx context.Context) error) {
		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()
		start := time.Now()
		err := call(ctx)
		elapsed := time.Since(start)
		mu.Lock()
		defer mu.Unlock()
		latencies = append(latencies, elapsed)
		if err != nil {
			failures++
		}
	}

	var wg sync.WaitGroup
	for tickAt := time.Now(); time.Now().Before(deadline); tickAt = tickAt.Add(tick) {
		// One GitHub call and one Linear call per tick.
		n := rand.IntN(1000) + *seed // deterministic-ish spread across ids
		wg.Add(2)
		go func() {
			defer wg.Done()
			run(func(ctx context.Context) error {
				_, err := github.GetIssue(ctx, &githubv1.GetIssueRequest{Owner: "acme", Repository: "widgets", IssueNumber: int64(n)})
				return err
			})
		}()
		go func() {
			defer wg.Done()
			run(func(ctx context.Context) error {
				_, err := linear.GetIssue(ctx, &linearv1.GetIssueRequest{Id: fmt.Sprintf("00000000-0000-4000-8000-%012d", n)})
				return err
			})
		}()
		time.Sleep(tick)
	}
	wg.Wait()

	lat := append([]time.Duration{}, latencies...)
	sort.Slice(lat, func(i, j int) bool { return lat[i] < lat[j] })
	pct := func(p float64) time.Duration {
		if len(lat) == 0 {
			return 0
		}
		return lat[int(float64(len(lat)-1)*p)]
	}
	n := len(latencies)
	errorRate := 0.0
	if n > 0 {
		errorRate = float64(failures) / float64(n)
	}
	fmt.Printf("requests=%d errors=%d error_rate=%.2f%% p50=%s p95=%s p99=%s max=%s\n",
		n, failures, errorRate*100, pct(0.50), pct(0.95), pct(0.99), pct(1.0))
}
