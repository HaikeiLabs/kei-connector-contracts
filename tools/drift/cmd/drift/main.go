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

// Command drift runs the nightly strict-unmarshal drift check against the
// connector sandboxes.
//
// For every connector in manifest/connector-manifest.v1.json it looks for
// sandbox credentials in KEI_SANDBOX_<CONNECTOR>_<FIELD> environment
// variables, a per-connector request file testenv/drift/<connector>.json
// naming safe READ operations, and a hand-written adapter that runs those
// operations through the generated Go ops. Every response body is
// re-unmarshalled with protojson DiscardUnknown=false: any unknown field is
// reported as drift (the contracts moved ahead of the sandbox or vice versa)
// and the process exits 1.
//
// Connectors without credentials, without a request file or without an
// adapter are skipped, never errors. Credential values are never printed;
// every output line is passed through a redaction pass over the credential
// values as a safety net.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	manifest := flag.String("manifest", "manifest/connector-manifest.v1.json", "path to connector-manifest.v1.json")
	requestsDir := flag.String("requests-dir", "testenv/drift", "directory of per-connector request files <connector>.json")
	only := flag.String("connector", "", "run only this connector id (default: all)")
	baseURL := flag.String("base-url", "", "override a connector base URL as connector=URL")
	timeout := flag.Duration("timeout", 60*time.Second, "per-operation timeout")
	flag.Parse()

	baseURLs := map[string]string{}
	if *baseURL != "" {
		id, url, ok := strings.Cut(*baseURL, "=")
		if !ok || id == "" || url == "" {
			fmt.Fprintln(os.Stderr, "drift: --base-url must be connector=URL")
			os.Exit(2)
		}
		baseURLs[id] = url
	}

	env := map[string]string{}
	for _, kv := range os.Environ() {
		k, v, _ := strings.Cut(kv, "=")
		env[k] = v
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	code := Run(ctx, Options{
		ManifestPath: *manifest,
		RequestsDir:  *requestsDir,
		Only:         *only,
		Env:          env,
		BaseURLs:     baseURLs,
		Out:          os.Stdout,
	})
	os.Exit(code)
}
