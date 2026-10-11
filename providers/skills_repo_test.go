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

package providers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

const skillsRepoTestCommit = "0123456789abcdef0123456789abcdef01234567"

func skillsRepoMeta(t *testing.T, config map[string]any) contract.Metadata {
	t.Helper()
	m := metaFor(t, contract.ProviderSkillsRepo, []string{"skills"}, nil)
	m.Config = config
	return m
}

func skillsRepoConfig() map[string]any {
	return map[string]any{"repo_url": "https://github.com/acme/agent-skills"}
}

func TestSkillsRepoSourceFromConfig(t *testing.T) {
	source, err := SkillsRepoSourceFromConfig(skillsRepoConfig())
	if err != nil {
		t.Fatal(err)
	}
	if source != (SkillsRepoSource{Host: "github.com", Owner: "acme", Repo: "agent-skills", Ref: "main", Subdirectory: "skills/"}) {
		t.Fatalf("defaults = %+v", source)
	}
	source, err = SkillsRepoSourceFromConfig(map[string]any{"repo_url": "https://github.com/acme/.github.git", "ref": "release/v1.2", "subdirectory": "/"})
	if err != nil || source.Repo != ".github" || source.Ref != "release/v1.2" || source.Subdirectory != "" || source.SkillPath("a") != "a" {
		t.Fatalf("root source = %+v, %v", source, err)
	}
	source, err = SkillsRepoSourceFromConfig(map[string]any{"repo_url": "https://github.com/acme/agent-skills", "subdirectory": "tools/skills"})
	if err != nil || source.Subdirectory != "tools/skills/" || source.SkillPath("a") != "tools/skills/a" {
		t.Fatalf("nested source = %+v, %v", source, err)
	}
	for name, config := range map[string]map[string]any{
		"missing url":       {},
		"non-string url":    {"repo_url": 7},
		"other host":        {"repo_url": "https://gitlab.com/acme/agent-skills"},
		"extra path":        {"repo_url": "https://github.com/acme/agent-skills/tree/main"},
		"dot-dot repo":      {"repo_url": "https://github.com/acme/.."},
		"query":             {"repo_url": "https://github.com/acme/agent-skills?ref=x"},
		"ref traversal":     {"repo_url": "https://github.com/acme/agent-skills", "ref": "main/../x"},
		"ref dot-dot":       {"repo_url": "https://github.com/acme/agent-skills", "ref": "a..b"},
		"ref option":        {"repo_url": "https://github.com/acme/agent-skills", "ref": "-x"},
		"ref reflog":        {"repo_url": "https://github.com/acme/agent-skills", "ref": "main@{1}"},
		"ref lock":          {"repo_url": "https://github.com/acme/agent-skills", "ref": "main.lock"},
		"ref empty":         {"repo_url": "https://github.com/acme/agent-skills", "ref": ""},
		"subdir traversal":  {"repo_url": "https://github.com/acme/agent-skills", "subdirectory": "skills/../.."},
		"subdir hidden":     {"repo_url": "https://github.com/acme/agent-skills", "subdirectory": ".git"},
		"subdir absolute":   {"repo_url": "https://github.com/acme/agent-skills", "subdirectory": "/skills"},
		"subdir backslash":  {"repo_url": "https://github.com/acme/agent-skills", "subdirectory": "skills\\x"},
		"subdir empty":      {"repo_url": "https://github.com/acme/agent-skills", "subdirectory": ""},
		"subdir non-string": {"repo_url": "https://github.com/acme/agent-skills", "subdirectory": true},
	} {
		if _, err := SkillsRepoSourceFromConfig(config); err == nil {
			t.Errorf("%s: config accepted", name)
		}
	}
}

func TestSkillsRepoResourceRejectsTraversal(t *testing.T) {
	if err := ValidateSkillsRepoResource("skills.list", "skills"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateSkillsRepoResource("skills.fetch", "skills/deploy_helper-2"); err != nil {
		t.Fatal(err)
	}
	for _, resource := range []string{"skills/..", "skills/.", "skills/%2e%2e", "skills/a/b", "skills/a/../b", "skills/.hidden", "skills/a.b", "skills/a\\b", "skills/", "skills", "repos/acme/x", "skills/" + strings.Repeat("a", 65)} {
		if err := ValidateSkillsRepoResource("skills.fetch", resource); err == nil {
			t.Errorf("resource %q accepted", resource)
		}
	}
	for _, name := range []string{"..", ".", "../etc", "a/b", "a\\b", "", "a.b"} {
		if err := ValidateSkillName(name); err == nil {
			t.Errorf("skill name %q accepted", name)
		}
	}

	client := NewSkillsRepo(MemorySkillsRepo{Archives: map[string][]byte{"..": []byte("x")}})
	meta := skillsRepoMeta(t, skillsRepoConfig())
	for _, resource := range []string{"skills/..", "skills/%2e%2e", "skills/a/../b"} {
		if _, err := client.Invoke(context.Background(), meta, invocation("skills.fetch", contract.ActionRead, resource), SkillsFetchPayload{}); err == nil {
			t.Errorf("fetch of %q was allowed", resource)
		}
	}
}

func TestSkillsRepoClientAllowlistsAndDigests(t *testing.T) {
	archive := []byte("archive-bytes")
	backend := MemorySkillsRepo{
		Commit:   skillsRepoTestCommit,
		Skills:   []SkillsRepoSkill{{Name: "deploy-helper", Version: "1.2.0", Description: "Deploys.", Harnesses: []string{"claude_code"}, Path: "skills/deploy-helper"}},
		Archives: map[string][]byte{"deploy-helper": archive},
	}
	client := NewSkillsRepo(backend)
	meta := skillsRepoMeta(t, map[string]any{"repo_url": "https://github.com/acme/agent-skills", "ref": "v1"})
	ctx := context.Background()

	result, err := client.Invoke(ctx, meta, invocation("skills.list", contract.ActionRead, "skills"), SkillsListPayload{})
	if err != nil {
		t.Fatal(err)
	}
	listing := result.Data.(SkillsRepoListing)
	if listing.Ref != "v1" || listing.Commit != skillsRepoTestCommit || len(listing.Skills) != 1 || listing.Skills[0].Name != "deploy-helper" {
		t.Fatalf("listing = %+v", listing)
	}
	encoded, _ := json.Marshal(listing.Skills[0])
	var fields map[string]any
	_ = json.Unmarshal(encoded, &fields)
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if strings.Join(keys, ",") != "description,harnesses,name,path,version" {
		t.Fatalf("listing fields = %v", keys)
	}

	result, err = client.Invoke(ctx, meta, invocation("skills.fetch", contract.ActionRead, "skills/deploy-helper"), SkillsFetchPayload{})
	if err != nil {
		t.Fatal(err)
	}
	got := result.Data.(SkillsRepoArchive)
	if got.Format != SkillsRepoArchiveFormat || got.Size != len(archive) || got.Digest != sha256Digest(archive) || !strings.HasPrefix(got.Digest, "sha256:") || got.Commit != skillsRepoTestCommit {
		t.Fatalf("archive = %+v", got)
	}

	if _, err := client.Invoke(ctx, meta, invocation("skills.fetch", contract.ActionRead, "skills/missing"), SkillsFetchPayload{}); err == nil {
		t.Error("missing skill was fetched")
	}
	if _, err := client.Invoke(ctx, meta, invocation("skills.fetch", contract.ActionRead, "skills/deploy-helper"), SkillsListPayload{}); err == nil {
		t.Error("payload mismatch was accepted")
	}
	if _, err := client.Invoke(ctx, meta, invocation("skills.list", contract.ActionRead, "skills/deploy-helper"), SkillsListPayload{}); err == nil {
		t.Error("list on a skill resource was accepted")
	}
	if _, err := client.Invoke(ctx, meta, invocation("skills.fetch", contract.ActionCreate, "skills/deploy-helper"), SkillsFetchPayload{}); err == nil {
		t.Error("fetch as a write was accepted")
	}
	if _, err := client.Invoke(ctx, skillsRepoMeta(t, nil), invocation("skills.list", contract.ActionRead, "skills"), SkillsListPayload{}); err == nil {
		t.Error("connector without a repository was accepted")
	}

	mismatched := NewSkillsRepo(renamingSkillsRepo{backend})
	if _, err := mismatched.Invoke(ctx, meta, invocation("skills.fetch", contract.ActionRead, "skills/deploy-helper"), SkillsFetchPayload{}); err == nil {
		t.Error("archive for another skill was accepted")
	}
}

type renamingSkillsRepo struct{ MemorySkillsRepo }

func (r renamingSkillsRepo) FetchSkill(ctx context.Context, source SkillsRepoSource, name string) (SkillsRepoArchive, error) {
	archive, err := r.MemorySkillsRepo.FetchSkill(ctx, source, name)
	archive.Skill.Name = "other"
	return archive, err
}

func TestParseSkillFrontmatter(t *testing.T) {
	fm, err := parseSkillFrontmatter([]byte("---\nname: deploy-helper\ndescription: \"Deploys: safely.\"\nversion: 1.2.0\nharnesses: [claude_code, 'codex']\nlicense: Apache-2.0\n---\n# Body\nname: ignored\n"))
	if err != nil || fm.Name != "deploy-helper" || fm.Description != "Deploys: safely." || fm.Version != "1.2.0" || strings.Join(fm.Harnesses, ",") != "claude_code,codex" {
		t.Fatalf("inline = %+v, %v", fm, err)
	}
	fm, err = parseSkillFrontmatter([]byte("---\nname: x\ndescription: >\n  Folded\n  text.\nallowed-tools:\n  - Bash\nmetadata:\n  version: \"2\"\n  harnesses:\n    - opencode\n    - pi\n---\n"))
	if err != nil || fm.Description != "Folded text." || fm.Version != "2" || strings.Join(fm.Harnesses, ",") != "opencode,pi" {
		t.Fatalf("nested = %+v, %v", fm, err)
	}
	for _, content := range []string{"", "# no frontmatter\n", "---\nname: x\n", "---\nname: " + strings.Repeat("a", SkillsRepoMaxSkillFileBytes) + "\n---\n"} {
		if _, err := parseSkillFrontmatter([]byte(content)); err == nil {
			t.Errorf("frontmatter %.20q accepted", content)
		}
	}

	skill, err := skillFromFrontmatter("x", "skills/x", []byte("---\nname: x\nversion: \"bad version\"\nharnesses: [claude_code, Bad, claude_code, \"a b\"]\ndescription: "+strings.Repeat("é", 2000)+"\n---\n"))
	if err != nil || skill.Version != "" || strings.Join(skill.Harnesses, ",") != "claude_code" || len([]rune(skill.Description)) != 1024 {
		t.Fatalf("sanitized = %+v, %v", skill, err)
	}
	if _, err := skillFromFrontmatter("x", "skills/x", []byte("---\nname: y\n---\n")); err == nil {
		t.Error("frontmatter naming another skill was accepted")
	}
}

// fakeSkillsGitHub serves the git data API for an in-memory repository at one
// commit. Object ids are content hashes, so the same tree always has the same
// id.
type fakeSkillsGitHub struct {
	t         *testing.T
	files     map[string]fakeSkillFile
	blobs     map[string]string
	trees     map[string]string
	root      string
	truncated bool
	requests  []string
}

type fakeSkillFile struct {
	content string
	mode    string
}

func newFakeSkillsGitHub(t *testing.T, files map[string]fakeSkillFile) *fakeSkillsGitHub {
	f := &fakeSkillsGitHub{t: t, files: files, blobs: map[string]string{}, trees: map[string]string{}}
	f.root = f.treeID("")
	for path, file := range files {
		f.blobs[objectID("blob:"+file.content)] = file.content
		for dir := parentDir(path); ; dir = parentDir(dir) {
			f.trees[f.treeID(dir)] = dir
			if dir == "" {
				break
			}
		}
	}
	return f
}

func objectID(s string) string {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func parentDir(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[:i]
	}
	return ""
}

func (f *fakeSkillsGitHub) treeID(dir string) string { return objectID("tree:" + dir) }

// entries lists dir's children, or every descendant when recursive, with
// paths relative to dir.
func (f *fakeSkillsGitHub) entries(dir string, recursive bool) []gitTreeEntry {
	prefix := ""
	if dir != "" {
		prefix = dir + "/"
	}
	seen := map[string]bool{}
	var out []gitTreeEntry
	for path, file := range f.files {
		if !strings.HasPrefix(path, prefix) {
			continue
		}
		rel := strings.TrimPrefix(path, prefix)
		segments := strings.Split(rel, "/")
		for i := 1; i < len(segments); i++ {
			sub := strings.Join(segments[:i], "/")
			if (recursive || i == 1) && !seen[sub] {
				seen[sub] = true
				out = append(out, gitTreeEntry{Path: sub, Mode: "040000", Type: "tree", SHA: f.treeID(prefix + sub)})
			}
		}
		if recursive || len(segments) == 1 {
			mode, kind := file.mode, "blob"
			if mode == "" {
				mode = "100644"
			}
			if mode == "160000" {
				kind = "commit"
			}
			out = append(out, gitTreeEntry{Path: rel, Mode: mode, Type: kind, SHA: objectID("blob:" + file.content), Size: int64(len(file.content))})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

func (f *fakeSkillsGitHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.requests = append(f.requests, r.URL.RequestURI())
	if r.Method != http.MethodGet || r.Header.Get("Authorization") != "Bearer fake-skills-token" || r.Header.Get("X-GitHub-Api-Version") != "2022-11-28" {
		f.t.Errorf("request = %s %s %v", r.Method, r.URL, r.Header)
	}
	const repo = "/repos/acme/agent-skills"
	switch path := r.URL.Path; {
	case path == repo+"/commits/main" || path == repo+"/commits/release/v1":
		_ = json.NewEncoder(w).Encode(map[string]any{"sha": skillsRepoTestCommit, "commit": map[string]any{"tree": map[string]any{"sha": f.root}}})
	case strings.HasPrefix(path, repo+"/git/trees/"):
		dir, ok := f.trees[strings.TrimPrefix(path, repo+"/git/trees/")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"truncated": f.truncated, "tree": f.entries(dir, r.URL.Query().Get("recursive") == "1")})
	case strings.HasPrefix(path, repo+"/git/blobs/"):
		if r.Header.Get("Accept") != "application/vnd.github.raw+json" {
			f.t.Errorf("blob accept = %q", r.Header.Get("Accept"))
		}
		content, ok := f.blobs[strings.TrimPrefix(path, repo+"/git/blobs/")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = io.WriteString(w, content)
	default:
		http.Error(w, `{"message":"Not Found fake-skills-token"}`, http.StatusNotFound)
	}
}

func skillsRepoFiles() map[string]fakeSkillFile {
	return map[string]fakeSkillFile{
		"README.md":                                {content: "repo readme"},
		"skills/deploy-helper/SKILL.md":            {content: "---\nname: deploy-helper\ndescription: Deploys the service.\nversion: 1.2.0\nharnesses: [claude_code, codex]\n---\nSteps.\n"},
		"skills/deploy-helper/scripts/run.sh":      {content: "#!/bin/sh\necho run\n", mode: "100755"},
		"skills/deploy-helper/references/guide.md": {content: "guide"},
		"skills/notes/readme.md":                   {content: "not a skill"},
		"skills/mismatch/SKILL.md":                 {content: "---\nname: other\n---\n"},
		"skills/Bad.Name/SKILL.md":                 {content: "---\nname: Bad.Name\n---\n"},
		"skills/reviewer/SKILL.md":                 {content: "---\nname: reviewer\ndescription: Reviews.\nmetadata:\n  harnesses:\n    - opencode\n---\n"},
		"skills/reviewer/nested/SKILL.md":          {content: "---\nname: nested\n---\n"},
	}
}

func skillsRepoHTTPClient(fake *fakeSkillsGitHub) (*SkillsRepoClient, *httptest.Server) {
	server := httptest.NewServer(fake)
	return NewSkillsRepoHTTP(testHTTPClient(server, func(t *testing.T, req *http.Request) {
		if req.URL.Scheme != "https" || req.URL.Host != "api.github.com" || req.Method != http.MethodGet {
			t.Errorf("request = %s %s", req.Method, req.URL)
		}
	}), &recordingResolver{token: "fake-skills-token"}), server
}

func TestGitHubSkillsRepoList(t *testing.T) {
	fake := newFakeSkillsGitHub(t, skillsRepoFiles())
	client, server := skillsRepoHTTPClient(fake)
	defer server.Close()

	result, err := client.Invoke(context.Background(), skillsRepoMeta(t, skillsRepoConfig()), invocation("skills.list", contract.ActionRead, "skills"), SkillsListPayload{})
	if err != nil {
		t.Fatal(err)
	}
	listing := result.Data.(SkillsRepoListing)
	want := []SkillsRepoSkill{
		{Name: "deploy-helper", Version: "1.2.0", Description: "Deploys the service.", Harnesses: []string{"claude_code", "codex"}, Path: "skills/deploy-helper"},
		{Name: "reviewer", Description: "Reviews.", Harnesses: []string{"opencode"}, Path: "skills/reviewer"},
	}
	gotJSON, _ := json.Marshal(listing.Skills)
	wantJSON, _ := json.Marshal(want)
	if listing.Commit != skillsRepoTestCommit || listing.Ref != "main" || string(gotJSON) != string(wantJSON) {
		t.Fatalf("listing = %s\nwant %s", gotJSON, wantJSON)
	}
	for _, request := range fake.requests {
		if strings.Contains(request, "README") || strings.Contains(request, "contents") {
			t.Errorf("unexpected request %q", request)
		}
	}
}

func TestGitHubSkillsRepoFetchBuildsDeterministicArchive(t *testing.T) {
	fake := newFakeSkillsGitHub(t, skillsRepoFiles())
	client, server := skillsRepoHTTPClient(fake)
	defer server.Close()
	meta := skillsRepoMeta(t, skillsRepoConfig())
	fetch := func() SkillsRepoArchive {
		t.Helper()
		result, err := client.Invoke(context.Background(), meta, invocation("skills.fetch", contract.ActionRead, "skills/deploy-helper"), SkillsFetchPayload{})
		if err != nil {
			t.Fatal(err)
		}
		return result.Data.(SkillsRepoArchive)
	}
	first, second := fetch(), fetch()
	if first.Digest != second.Digest || first.Digest != sha256Digest(first.Archive) || first.Size != len(first.Archive) {
		t.Fatalf("digests = %s, %s", first.Digest, second.Digest)
	}
	if first.Skill.Name != "deploy-helper" || first.Skill.Version != "1.2.0" || first.Commit != skillsRepoTestCommit || first.Format != "tar+gzip" {
		t.Fatalf("archive = %+v", first.Skill)
	}

	gz, err := gzip.NewReader(bytes.NewReader(first.Archive))
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gz)
	got := map[string]string{}
	modes := map[string]int64{}
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		content, _ := io.ReadAll(tr)
		got[header.Name], modes[header.Name] = string(content), header.Mode
	}
	files := skillsRepoFiles()
	want := map[string]string{
		"deploy-helper/SKILL.md":            files["skills/deploy-helper/SKILL.md"].content,
		"deploy-helper/scripts/run.sh":      files["skills/deploy-helper/scripts/run.sh"].content,
		"deploy-helper/references/guide.md": "guide",
	}
	if len(got) != len(want) {
		t.Fatalf("archive entries = %v", got)
	}
	for name, content := range want {
		if got[name] != content {
			t.Errorf("%s = %q, want %q", name, got[name], content)
		}
	}
	if modes["deploy-helper/scripts/run.sh"] != 0o755 || modes["deploy-helper/SKILL.md"] != 0o644 {
		t.Errorf("modes = %v", modes)
	}
}

func TestGitHubSkillsRepoNestedSubdirectoryAndRef(t *testing.T) {
	fake := newFakeSkillsGitHub(t, map[string]fakeSkillFile{
		"tools/skills/x/SKILL.md": {content: "---\nname: x\n---\n"},
	})
	client, server := skillsRepoHTTPClient(fake)
	defer server.Close()
	meta := skillsRepoMeta(t, map[string]any{"repo_url": "https://github.com/acme/agent-skills", "ref": "release/v1", "subdirectory": "tools/skills/"})
	result, err := client.Invoke(context.Background(), meta, invocation("skills.list", contract.ActionRead, "skills"), SkillsListPayload{})
	if err != nil {
		t.Fatal(err)
	}
	listing := result.Data.(SkillsRepoListing)
	if len(listing.Skills) != 1 || listing.Skills[0].Path != "tools/skills/x" || listing.Ref != "release/v1" {
		t.Fatalf("listing = %+v", listing)
	}
}

func TestGitHubSkillsRepoFailsClosed(t *testing.T) {
	meta := skillsRepoMeta(t, skillsRepoConfig())
	ctx := context.Background()
	for name, tc := range map[string]struct {
		files     map[string]fakeSkillFile
		truncated bool
		skill     string
		want      string
	}{
		"symlink":        {files: map[string]fakeSkillFile{"skills/x/SKILL.md": {content: "---\nname: x\n---\n"}, "skills/x/link": {content: "../../etc", mode: "120000"}}, skill: "x", want: "symlink"},
		"submodule":      {files: map[string]fakeSkillFile{"skills/x/SKILL.md": {content: "---\nname: x\n---\n"}, "skills/x/vendor": {content: "sub", mode: "160000"}}, skill: "x", want: "submodule"},
		"no SKILL.md":    {files: map[string]fakeSkillFile{"skills/x/readme.md": {content: "x"}}, skill: "x", want: "SKILL.md"},
		"name mismatch":  {files: map[string]fakeSkillFile{"skills/x/SKILL.md": {content: "---\nname: y\n---\n"}}, skill: "x", want: "does not match"},
		"missing folder": {files: map[string]fakeSkillFile{"skills/x/SKILL.md": {content: "---\nname: x\n---\n"}}, skill: "y", want: "not found"},
		"truncated tree": {files: map[string]fakeSkillFile{"skills/x/SKILL.md": {content: "---\nname: x\n---\n"}}, truncated: true, skill: "x", want: "too large"},
		"oversized file": {files: map[string]fakeSkillFile{"skills/x/SKILL.md": {content: "---\nname: x\n---\n"}, "skills/x/big.bin": {content: strings.Repeat("a", SkillsRepoMaxFileBytes+1)}}, skill: "x", want: "larger"},
	} {
		fake := newFakeSkillsGitHub(t, tc.files)
		fake.truncated = tc.truncated
		client, server := skillsRepoHTTPClient(fake)
		_, err := client.Invoke(ctx, meta, invocation("skills.fetch", contract.ActionRead, "skills/"+tc.skill), SkillsFetchPayload{})
		server.Close()
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: error = %v, want %q", name, err, tc.want)
		}
	}

	for _, path := range []string{"../x", "a//b", "./a", "a\\b"} {
		if _, err := skillFiles([]gitTreeEntry{{Path: path, Type: "blob", SHA: skillsRepoTestCommit}}); err == nil {
			t.Errorf("tree path %q accepted", path)
		}
	}
}

func TestGitHubSkillsRepoRedactsErrorsAndNeedsCredential(t *testing.T) {
	meta := skillsRepoMeta(t, skillsRepoConfig())
	ctx := context.Background()
	client, err := New(contract.ProviderSkillsRepo)
	if err != nil || client.Provider() != contract.ProviderSkillsRepo {
		t.Fatalf("New = %v, %v", client, err)
	}
	if _, err := client.Invoke(ctx, meta, invocation("skills.list", contract.ActionRead, "skills"), SkillsListPayload{}); err == nil || !strings.Contains(err.Error(), "credential") {
		t.Fatalf("unconfigured resolver error = %v", err)
	}

	noNetwork := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("must not send") })}
	if _, err := NewSkillsRepoHTTP(noNetwork, credentialErrorResolver{}).Invoke(ctx, meta, invocation("skills.list", contract.ActionRead, "skills"), SkillsListPayload{}); err == nil || err.Error() != "credential resolution failed" {
		t.Fatalf("resolver error = %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message":"Bad credentials fake-skills-token"}`, http.StatusUnauthorized)
	}))
	defer server.Close()
	rejected := NewSkillsRepoHTTP(testHTTPClient(server, func(*testing.T, *http.Request) {}), &recordingResolver{token: "fake-skills-token"})
	for _, tc := range []struct {
		capability, resource string
		payload              Payload
	}{{"skills.list", "skills", SkillsListPayload{}}, {"skills.fetch", "skills/x", SkillsFetchPayload{}}} {
		_, err := rejected.Invoke(ctx, meta, invocation(tc.capability, contract.ActionRead, tc.resource), tc.payload)
		if err == nil || strings.Contains(err.Error(), "fake-skills-token") || strings.Contains(err.Error(), "Bad credentials") {
			t.Errorf("%s rejection error = %v", tc.capability, err)
		}
	}
}
