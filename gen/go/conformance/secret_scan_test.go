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

// Secret-safety conformance (HAI-1001, plan task 1.4).
//
// The contract repository is public and metadata only: no credential value
// may ever be checked in. TestNoCredentialMaterialInContracts walks proto/,
// fixtures/ and gen/ and fails on token-shaped strings. The planted-token
// tests prove the scanner actually detects them: they build every token
// shape at test time (no token-shaped literal is committed) and feed it
// through the scanner, one shape at a time and end to end through a temp
// file.

package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// contractTrees are the trees the secret scan covers, relative to the
// repository root: the hand-written proto sources, the fixtures, and the
// generated bindings. The root module's own packages and their test data are
// not contract material and are out of scope.
func contractTrees(t *testing.T) map[string]string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	return map[string]string{
		"proto":    filepath.Join(root, "proto"),
		"fixtures": filepath.Join(root, "fixtures"),
		"gen":      filepath.Join(root, "gen"),
	}
}

// TestNoCredentialMaterialInContracts walks the contract trees and fails on
// any token-shaped string. Findings name the tree, file, line and pattern
// only; the scanner never echoes the matched text (the same rule
// kei-connector-runtime internal/secrets/ref.go applies to error messages).
func TestNoCredentialMaterialInContracts(t *testing.T) {
	for tree, root := range contractTrees(t) {
		findings, err := scanTree(tree, root)
		if err != nil {
			t.Fatalf("scanning %s: %v", tree, err)
		}
		for _, f := range findings {
			t.Errorf("credential-shaped string in %s tree: %s", tree, f)
		}
	}
}

// plantedToken is one synthetic token shape, built at test time from
// concatenated fragments so no token-shaped literal enters this file. want
// names the scanner pattern that must fire for it.
type plantedToken struct {
	name   string
	build  func() string
	want   string
	skipOn bool // detected only via the quoted high-entropy path, not a prefix shape
}

// repeat builds an n-times-repeated fragment at runtime.
func repeat(s string, n int) string { return strings.Repeat(s, n) }

// highEntropyPlantedToken builds a 64-character mixed-alphabet token at
// runtime. A literal 40+ character high-entropy token must not appear in
// this test source, so the token is assembled by stepping through the
// alphabet.
func highEntropyPlantedToken() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	for i := 0; i < 64; i++ {
		b.WriteByte(alphabet[(i*11)%len(alphabet)])
	}
	return b.String()
}

var plantedTokens = []plantedToken{
	{"github classic token", func() string { return "ghp_" + repeat("x3F9kQm2Lp", 4) }, "github-token", false},
	{"github fine-grained pat", func() string { return "github_pat_" + repeat("Ab12Cd34", 5) }, "github-fine-grained-pat", false},
	{"linear api key", func() string { return "lin_api_" + repeat("7f3a9c1d", 5) }, "linear-api-key", false},
	{"slack bot token", func() string { return "xoxb-" + repeat("1234567890", 5) }, "slack-token", false},
	{"slack refresh token", func() string { return "xoxr-" + repeat("1a2b3c4d5e", 3) }, "slack-token", false},
	{"slack app token", func() string { return "xapp-" + repeat("1a2b3c4d5e", 4) }, "slack-app-token", false},
	{"openai-style key", func() string { return "sk-" + repeat("Qw7Np2Lm5X", 6) }, "openai-style-key", false},
	{"stripe live key", func() string { return "sk_live_" + repeat("Kj8mN2Qp", 5) }, "stripe-style-key", false},
	{"stripe publishable key", func() string { return "pk_test_" + repeat("Yh4tR9Wq", 5) }, "stripe-style-key", false},
	{"aws access key id", func() string { return "AKIA" + repeat("J4X7Q2M9", 2) }, "aws-access-key-id", false},
	{"google api key", func() string { return "AIza" + repeat("aB3xY9Kq", 5) }, "google-api-key", false},
	// PEM header, assembled at runtime so the full shape never appears in
	// this source file (the scanner scans gen/ too).
	{"pem private key header", func() string { return "-----BEGIN " + "RSA " + "PRIVATE " + "KEY" + "-----" }, "pem-private-key", false},
	{"compact jwt", func() string {
		head := "eyJ" + repeat("hbGciOi", 4)
		body := repeat("c2VjcmV0", 5)
		tail := repeat("dG9rZW4", 6)
		return head + "." + body + "." + tail
	}, "jwt", false},
	{"url with embedded credentials", func() string {
		return "postgres://" + "svcuser" + ":" + "hunter2" + "@" + "db.example.invalid/crm"
	}, "url-embedded-credentials", false},
	{"high-entropy token", highEntropyPlantedToken, "high-entropy-token", true},
}

// benignLines are strings the scanner must not flag: ordinary URLs, synthetic
// ids, and connection references without inline credentials.
var benignLines = []string{
	`"html_url": "https://example.invalid/example-org/example-repo/issues/1347"`,
	`"node_id": "synthetic-node-id-1"`,
	"postgres://localhost:5432/contracts",
	"https://api.example.invalid/v1",
	"AKIA is the AWS access key id prefix",
	`{"etag": "synthetic-etag-1"}`,
}

func TestPlantedSecretMaterialIsDetected(t *testing.T) {
	for _, tt := range plantedTokens {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.build()
			// The high-entropy catch-all matches quoted string values, the
			// way a leaked secret appears in a JSON payload; feed it that way.
			feed := token + "\n"
			if tt.skipOn {
				feed = `{"api_key": "` + token + `"}` + "\n"
			}
			findings := scanText(feed, false)
			if len(findings) == 0 {
				t.Fatalf("scanner detected nothing for the planted %s", tt.name)
			}
			if !strings.HasSuffix(findings[0], tt.want) {
				t.Fatalf("planted %s flagged as %v, want pattern %s", tt.name, findings, tt.want)
			}
			for _, f := range findings {
				if strings.Contains(f, token) {
					t.Errorf("finding echoes the planted material: %s", f)
				}
			}
		})
	}
}

// TestPlantedFileIsDetectedEndToEnd writes one planted file into a temp
// directory (never the repository) and walks it with the same tree scanner
// the repository test uses, proving the walk + scan pipeline end to end.
func TestPlantedFileIsDetectedEndToEnd(t *testing.T) {
	var b strings.Builder
	for _, tt := range plantedTokens {
		if tt.skipOn {
			b.WriteString(`{"api_key": "` + tt.build() + `"}` + "\n")
			continue
		}
		b.WriteString(tt.build() + "\n")
	}
	for _, line := range benignLines {
		b.WriteString(line + "\n")
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "planted.json"), []byte(b.String()), 0o600); err != nil {
		t.Fatalf("writing planted file: %v", err)
	}

	findings, err := scanTree("planted", dir)
	if err != nil {
		t.Fatalf("scanning planted tree: %v", err)
	}

	found := make(map[string]bool, len(findings))
	for _, f := range findings {
		pattern := f[strings.LastIndex(f, " ")+1:]
		found[pattern] = true
		if strings.Contains(f, "hunter2") {
			t.Errorf("finding echoes the planted material: %s", f)
		}
	}
	for _, tt := range plantedTokens {
		if !found[tt.want] {
			t.Errorf("planted %s was not detected (want pattern %s); findings: %v", tt.name, tt.want, findings)
		}
	}
}

// TestBenignContentIsNotFlagged keeps the scanner precise: ordinary URLs,
// synthetic ids and references without credentials must produce nothing.
func TestBenignContentIsNotFlagged(t *testing.T) {
	text := strings.Join(benignLines, "\n") + "\n"
	if findings := scanText(text, false); len(findings) != 0 {
		t.Fatalf("benign content flagged: %v", findings)
	}
}
