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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

var gitObjectID = regexp.MustCompile(`^([0-9a-f]{40}|[0-9a-f]{64})$`)

const (
	gitModeExecutable = "100755"
	gitModeSymlink    = "120000"
)

// GitHubSkillsRepo is the GitHub backend for a skills repository. It reads
// through the fixed GitHub API base with GET requests only: it resolves the
// ref to a commit, walks the git trees to the skills folder, and reads blobs
// by object id, so every read in one call comes from the same commit. The
// token is the invoking user's OAuth token (per_user) or the connector's
// shared-secret token, resolved per invocation and never returned.
type GitHubSkillsRepo struct {
	httpClient  *http.Client
	credentials CredentialResolver
}

func NewAuthenticatedSkillsRepo(config RuntimeConfig) *SkillsRepoClient {
	config = normalizeRuntimeConfig(config)
	return NewSkillsRepo(&GitHubSkillsRepo{httpClient: config.HTTPClient, credentials: config.Credentials})
}

func NewSkillsRepoHTTP(httpClient *http.Client, credentials CredentialResolver) *SkillsRepoClient {
	return NewAuthenticatedSkillsRepo(RuntimeConfig{HTTPClient: httpClient, Credentials: credentials})
}

type gitTreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	Size int64  `json:"size"`
}

// skillsRepoRead is one invocation's reads: one token, one repository.
type skillsRepoRead struct {
	backend *GitHubSkillsRepo
	token   string
	source  SkillsRepoSource
}

func (b *GitHubSkillsRepo) open(ctx context.Context, source SkillsRepoSource) (*skillsRepoRead, error) {
	if source.Host != skillsRepoGitHubHost {
		return nil, errors.New("skills repository host is not supported")
	}
	scope, err := scopeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	token, err := b.credentials.Resolve(ctx, scope.tenantID, scope.workspaceID, scope.subject, scope.credential)
	if err != nil || token == "" {
		return nil, errors.New("credential resolution failed")
	}
	return &skillsRepoRead{backend: b, token: token, source: source}, nil
}

func (b *GitHubSkillsRepo) ListSkills(ctx context.Context, source SkillsRepoSource) (SkillsRepoListing, error) {
	read, err := b.open(ctx, source)
	if err != nil {
		return SkillsRepoListing{}, err
	}
	commit, root, err := read.resolve(ctx)
	if err != nil {
		return SkillsRepoListing{}, err
	}
	folder, err := read.walk(ctx, root, strings.Split(strings.TrimSuffix(source.Subdirectory, "/"), "/"))
	if err != nil {
		return SkillsRepoListing{}, err
	}
	entries, err := read.tree(ctx, folder, true)
	if err != nil {
		return SkillsRepoListing{}, err
	}
	// A skill is a folder directly under the skills folder with a SKILL.md.
	// Other folders are not skills and are left out.
	var skillMDs []gitTreeEntry
	for _, entry := range entries {
		name, file, ok := strings.Cut(entry.Path, "/")
		if ok && file == "SKILL.md" && skillName.MatchString(name) {
			skillMDs = append(skillMDs, entry)
		}
	}
	if len(skillMDs) > SkillsRepoMaxSkills {
		return SkillsRepoListing{}, fmt.Errorf("skills repository has more than %d skills", SkillsRepoMaxSkills)
	}
	sort.Slice(skillMDs, func(i, j int) bool { return skillMDs[i].Path < skillMDs[j].Path })
	skills := []SkillsRepoSkill{}
	for _, entry := range skillMDs {
		name := strings.TrimSuffix(entry.Path, "/SKILL.md")
		entry.Path = "SKILL.md"
		file, ok := findSkillFile([]gitTreeEntry{entry})
		if !ok {
			continue
		}
		content, err := read.blob(ctx, file)
		if err != nil {
			return SkillsRepoListing{}, err
		}
		// A SKILL.md without valid frontmatter, or naming another skill, is
		// left out of the listing; skills.fetch refuses it.
		skill, err := skillFromFrontmatter(name, source.SkillPath(name), content)
		if err != nil {
			continue
		}
		skills = append(skills, skill)
	}
	return SkillsRepoListing{Ref: source.Ref, Commit: commit, Skills: skills}, nil
}

func (b *GitHubSkillsRepo) FetchSkill(ctx context.Context, source SkillsRepoSource, name string) (SkillsRepoArchive, error) {
	if err := ValidateSkillName(name); err != nil {
		return SkillsRepoArchive{}, err
	}
	read, err := b.open(ctx, source)
	if err != nil {
		return SkillsRepoArchive{}, err
	}
	commit, root, err := read.resolve(ctx)
	if err != nil {
		return SkillsRepoArchive{}, err
	}
	segments := append(strings.Split(strings.TrimSuffix(source.Subdirectory, "/"), "/"), name)
	folder, err := read.walk(ctx, root, segments)
	if err != nil {
		return SkillsRepoArchive{}, err
	}
	entries, err := read.tree(ctx, folder, true)
	if err != nil {
		return SkillsRepoArchive{}, err
	}
	files, err := skillFiles(entries)
	if err != nil {
		return SkillsRepoArchive{}, err
	}
	skillFile, ok := findSkillFile(files)
	if !ok {
		return SkillsRepoArchive{}, errors.New("skill has no SKILL.md")
	}
	content, err := read.blob(ctx, skillFile)
	if err != nil {
		return SkillsRepoArchive{}, err
	}
	skill, err := skillFromFrontmatter(name, source.SkillPath(name), content)
	if err != nil {
		return SkillsRepoArchive{}, err
	}
	archive, err := read.archive(ctx, name, files)
	if err != nil {
		return SkillsRepoArchive{}, err
	}
	return SkillsRepoArchive{Skill: skill, Commit: commit, Format: SkillsRepoArchiveFormat, Archive: archive}, nil
}

// skillFiles checks a skill folder's recursive tree and returns its files in
// path order. Symlinks and submodules are refused, because either can point
// outside the folder; so are paths that are not plain relative paths, and
// folders past the file and size bounds.
func skillFiles(entries []gitTreeEntry) ([]gitTreeEntry, error) {
	var files []gitTreeEntry
	var total int64
	for _, entry := range entries {
		if !plainRelativePath(entry.Path) {
			return nil, errors.New("skill contains an invalid path")
		}
		switch {
		case entry.Type == "tree":
			continue
		case entry.Type == "commit" || entry.Mode == gitModeSymlink:
			return nil, errors.New("skill contains a symlink or submodule")
		case entry.Type != "blob":
			return nil, errors.New("skill contains an unsupported entry")
		}
		if entry.Size < 0 || entry.Size > SkillsRepoMaxFileBytes {
			return nil, fmt.Errorf("skill file is larger than %d bytes", SkillsRepoMaxFileBytes)
		}
		total += entry.Size
		files = append(files, entry)
	}
	if len(files) > SkillsRepoMaxFiles {
		return nil, fmt.Errorf("skill has more than %d files", SkillsRepoMaxFiles)
	}
	if total > SkillsRepoMaxArchiveBytes {
		return nil, fmt.Errorf("skill is larger than %d bytes", SkillsRepoMaxArchiveBytes)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func plainRelativePath(path string) bool {
	if path == "" || strings.ContainsAny(path, "\\\x00") {
		return false
	}
	for _, segment := range strings.Split(path, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
	}
	return true
}

func findSkillFile(entries []gitTreeEntry) (gitTreeEntry, bool) {
	for _, entry := range entries {
		if entry.Path == "SKILL.md" && entry.Type == "blob" && entry.Mode != gitModeSymlink && entry.Size <= SkillsRepoMaxSkillFileBytes {
			return entry, true
		}
	}
	return gitTreeEntry{}, false
}

// archive writes files as a deterministic tar+gzip: path order, fixed
// timestamps and owners, and modes 0644 or 0755, so one commit always yields
// the same digest.
func (r *skillsRepoRead) archive(ctx context.Context, name string, files []gitTreeEntry) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, file := range files {
		content, err := r.blob(ctx, file)
		if err != nil {
			return nil, err
		}
		mode := int64(0o644)
		if file.Mode == gitModeExecutable {
			mode = 0o755
		}
		header := &tar.Header{Typeflag: tar.TypeReg, Name: name + "/" + file.Path, Mode: mode, Size: int64(len(content)), ModTime: time.Unix(0, 0), Format: tar.FormatPAX}
		if err := tw.WriteHeader(header); err != nil {
			return nil, errors.New("skill archive could not be written")
		}
		if _, err := tw.Write(content); err != nil {
			return nil, errors.New("skill archive could not be written")
		}
	}
	if err := tw.Close(); err != nil {
		return nil, errors.New("skill archive could not be written")
	}
	if err := gz.Close(); err != nil {
		return nil, errors.New("skill archive could not be written")
	}
	return buf.Bytes(), nil
}

// resolve returns the commit the source ref names and its root tree.
func (r *skillsRepoRead) resolve(ctx context.Context) (commit, tree string, err error) {
	var response struct {
		SHA    string `json:"sha"`
		Commit struct {
			Tree struct {
				SHA string `json:"sha"`
			} `json:"tree"`
		} `json:"commit"`
	}
	refPath := make([]string, 0, 4)
	for _, segment := range strings.Split(r.source.Ref, "/") {
		refPath = append(refPath, url.PathEscape(segment))
	}
	if err := r.getJSON(ctx, r.repoPath()+"/commits/"+strings.Join(refPath, "/"), &response); err != nil {
		return "", "", err
	}
	if !gitObjectID.MatchString(response.SHA) || !gitObjectID.MatchString(response.Commit.Tree.SHA) {
		return "", "", errors.New("provider response was invalid")
	}
	return response.SHA, response.Commit.Tree.SHA, nil
}

// walk follows segments from the tree root and returns the tree they name.
// An empty first segment is the root itself.
func (r *skillsRepoRead) walk(ctx context.Context, root string, segments []string) (string, error) {
	current := root
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		entries, err := r.tree(ctx, current, false)
		if err != nil {
			return "", err
		}
		next := ""
		for _, entry := range entries {
			if entry.Path == segment && entry.Type == "tree" {
				next = entry.SHA
			}
		}
		if next == "" {
			return "", errors.New("skills repository path not found")
		}
		current = next
	}
	return current, nil
}

func (r *skillsRepoRead) tree(ctx context.Context, sha string, recursive bool) ([]gitTreeEntry, error) {
	if !gitObjectID.MatchString(sha) {
		return nil, errors.New("provider response was invalid")
	}
	path := r.repoPath() + "/git/trees/" + sha
	if recursive {
		path += "?recursive=1"
	}
	var response struct {
		Truncated bool           `json:"truncated"`
		Tree      []gitTreeEntry `json:"tree"`
	}
	if err := r.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}
	if response.Truncated {
		return nil, errors.New("skills repository tree is too large")
	}
	for _, entry := range response.Tree {
		if !gitObjectID.MatchString(entry.SHA) {
			return nil, errors.New("provider response was invalid")
		}
	}
	return response.Tree, nil
}

// blob reads one file's raw content and checks it against the tree's size.
func (r *skillsRepoRead) blob(ctx context.Context, entry gitTreeEntry) ([]byte, error) {
	if !gitObjectID.MatchString(entry.SHA) {
		return nil, errors.New("provider response was invalid")
	}
	body, err := r.get(ctx, r.repoPath()+"/git/blobs/"+entry.SHA, "application/vnd.github.raw+json", SkillsRepoMaxFileBytes)
	if err != nil {
		return nil, err
	}
	if int64(len(body)) != entry.Size {
		return nil, errors.New("provider response was invalid")
	}
	return body, nil
}

func (r *skillsRepoRead) repoPath() string {
	return "/repos/" + url.PathEscape(r.source.Owner) + "/" + url.PathEscape(r.source.Repo)
}

func (r *skillsRepoRead) getJSON(ctx context.Context, path string, out any) error {
	body, err := r.get(ctx, path, "application/vnd.github+json", SkillsRepoMaxArchiveBytes)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, out); err != nil {
		return errors.New("provider response was invalid")
	}
	return nil
}

// get sends one GET to the GitHub API and reads at most limit bytes. Errors
// never carry the response body or the token.
func (r *skillsRepoRead) get(ctx context.Context, path, accept string, limit int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubAPIBase+path, nil)
	if err != nil {
		return nil, errors.New("provider request could not be created")
	}
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Accept", accept)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := r.backend.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("provider request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("skills repository path not found")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("provider request was rejected")
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, errors.New("provider response was invalid")
	}
	if int64(len(body)) > limit {
		return nil, errors.New("provider response is too large")
	}
	return body, nil
}

func sha256Digest(content []byte) string {
	sum := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(sum[:])
}

var _ SkillsRepoBackend = (*GitHubSkillsRepo)(nil)
