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

package conformance

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

// scanPattern is one deterministic credential fingerprint: a token shape
// that only credential material produces.
type scanPattern struct {
	name string
	re   *regexp.Regexp
}

// scanPatterns mirrors the secret-material heuristics of kei-connector-runtime
// internal/secrets/ref.go (knownSecretPrefixes, jwtShape, urlWithEmbeddedCreds)
// adapted from whole-reference prefix checks to in-text token shapes: each
// pattern requires a plausible token body after the fingerprint so ordinary
// prose or identifiers never match.
var scanPatterns = []scanPattern{
	// GitHub personal access tokens: classic ghp_/gho_/ghu_/ghs_/ghr_ and
	// fine-grained github_pat_.
	{"github-token", regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{20,}`)},
	{"github-fine-grained-pat", regexp.MustCompile(`github_pat_[A-Za-z0-9_]{20,}`)},
	// Linear API keys.
	{"linear-api-key", regexp.MustCompile(`lin_api_[A-Za-z0-9]{20,}`)},
	// Slack bot/user tokens and refresh tokens (ref.go: xoxb-, xoxp-, xoxa-).
	{"slack-token", regexp.MustCompile(`xox[abpr]-[A-Za-z0-9-]{10,}`)},
	// Slack app-level tokens (ref.go: xapp-).
	{"slack-app-token", regexp.MustCompile(`xapp-[A-Za-z0-9-]{10,}`)},
	// OpenAI-style keys (ref.go prefix sk-).
	{"openai-style-key", regexp.MustCompile(`sk-[A-Za-z0-9]{20,}`)},
	// Stripe-style live/test secret and publishable keys
	// (ref.go: sk_live_, sk_test_, pk_live_, pk_test_).
	{"stripe-style-key", regexp.MustCompile(`(?:sk|pk)_(?:live|test)_[A-Za-z0-9]{16,}`)},
	// AWS access key ids (ref.go: AKIA, ASIA).
	{"aws-access-key-id", regexp.MustCompile(`(?:AKIA|ASIA)[A-Z0-9]{16}`)},
	// Google API keys (ref.go: AIza).
	{"google-api-key", regexp.MustCompile(`AIza[0-9A-Za-z_-]{35}`)},
	// PEM private key blocks (ref.go prefix -----BEGIN).
	{"pem-private-key", regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
	// Compact JWT: three base64url segments with an eyJ ({") header
	// (ref.go jwtShape).
	{"jwt", regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)},
	// URL whose userinfo carries inline credentials (ref.go
	// urlWithEmbeddedCreds: a URL embedding a user and password).
	{"url-embedded-credentials", regexp.MustCompile(`://[^/@\s]+:[^/@\s]+@`)},
}

// highEntropyQuoted matches a double-quoted string value made only of token
// characters (the base64url alphabet plus hex and underscores). The entropy
// filter decides whether the value is credential material.
var highEntropyQuoted = regexp.MustCompile(`"([A-Za-z0-9_-]{40,})"`)

const (
	entropyMinLen   = 40
	entropyMinScore = 3.5
)

// highEntropyToken flags a delimited token with the entropy profile of an
// encoded secret (base64/hex), the same bar ref.go applies to a whole
// reference: at least 40 characters, no structural separators (/ : .), and
// Shannon entropy of 3.5 bits per character or more.
func highEntropyToken(s string) bool {
	if len(s) < entropyMinLen {
		return false
	}
	if strings.ContainsAny(s, "/:.") {
		return false
	}
	return shannonEntropy(s) >= entropyMinScore
}

// shannonEntropy is the entropy of s in bits per character, ported from
// kei-connector-runtime internal/secrets/ref.go.
func shannonEntropy(s string) float64 {
	if s == "" {
		return 0
	}
	freq := make(map[byte]float64, len(s))
	for i := 0; i < len(s); i++ {
		freq[s[i]]++
	}
	total := float64(len(s))
	entropy := 0.0
	for _, count := range freq {
		p := count / total
		entropy -= p * math.Log2(p)
	}
	return entropy
}

// scanText scans one text document line by line. Findings name the line and
// the pattern only — never the matched text. skipEntropy disables the
// high-entropy catch-all for the document (see entropySkip).
func scanText(text string, skipEntropy bool) []string {
	var findings []string
	seen := make(map[string]bool)
	add := func(f string) {
		if !seen[f] {
			seen[f] = true
			findings = append(findings, f)
		}
	}
	for i, line := range strings.Split(text, "\n") {
		for _, p := range scanPatterns {
			if p.re.MatchString(line) {
				add(fmt.Sprintf("line %d: %s", i+1, p.name))
			}
		}
		if skipEntropy {
			continue
		}
		for _, m := range highEntropyQuoted.FindAllStringSubmatch(line, -1) {
			if highEntropyToken(m[1]) {
				add(fmt.Sprintf("line %d: high-entropy-token", i+1))
			}
		}
	}
	return findings
}

// scanTree walks root and scans every text file, returning findings as
// "label/rel: line N: pattern". The label names the tree being walked.
func scanTree(label, root string) ([]string, error) {
	var findings []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "__pycache__", "node_modules", ".git", "target":
				return filepath.SkipDir
			}
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !utf8.Valid(data) {
			return nil // binary artifact (pyc, images)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		for _, f := range scanText(string(data), entropySkip(rel)) {
			findings = append(findings, fmt.Sprintf("%s/%s: %s", label, rel, f))
		}
		return nil
	})
	return findings, err
}

// entropySkip reports whether the high-entropy catch-all is disabled for a
// file. Derived files that provably embed long encoded blobs that are not
// credentials are skipped: serialized FileDescriptorProto base64 in generated
// bindings (derived from proto/, which is scanned directly) and dependency
// integrity checksums of public packages in lockfiles.
func entropySkip(rel string) bool {
	base := filepath.Base(rel)
	if base == "go.sum" || base == "Cargo.lock" {
		return true
	}
	return strings.HasSuffix(rel, "_pb2.py") || strings.HasSuffix(rel, "_pb.ts")
}
