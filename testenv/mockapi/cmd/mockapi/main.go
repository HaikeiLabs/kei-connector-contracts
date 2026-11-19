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

// Command mockapi serves the contract-derived provider mock on an address.
//
//	mockapi -addr :8080 -manifest ../../manifest/connector-manifest.v1.json -fixtures ../../fixtures/operations
package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/HaikeiLabs/kei-connector-contracts/testenv/mockapi"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	manifest := flag.String("manifest", "../../manifest/connector-manifest.v1.json", "connector manifest path")
	fixtures := flag.String("fixtures", "../../fixtures/operations", "operations fixtures directory")
	seed := flag.Int("seed", 1, "seed for the synthetic page ids")
	flag.Parse()
	srv, err := mockapi.New(mockapi.Options{ManifestPath: *manifest, FixturesDir: *fixtures, Seed: *seed})
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	// The actual bound address goes to stderr: -addr 127.0.0.1:0 picks a
	// free port for the language smoke tests.
	log.Printf("mockapi: serving %d routes on %s", len(srv.Routes()), listener.Addr())
	log.Fatal(http.Serve(listener, srv))
}
