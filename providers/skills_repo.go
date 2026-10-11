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
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

// A skills repository is a git repository holding one folder per agent skill
// under a skills folder, each with a SKILL.md whose frontmatter names the
// skill. The connector's Config is its source (see SkillsRepoSourceFromConfig).
//
// Resources:
//
//	skills.list   skills
//	skills.fetch  skills/<name>
//
// A skill name is one path segment of letters, digits, '-' and '_'. It can
// never contain '.', '/' or '\', so a fetch cannot leave the skills folder.

const (
	// SkillsRepoArchiveFormat is the only skills.fetch archive format: a gzip
	// compressed tar whose entries are <name>/<path inside the skill folder>.
	SkillsRepoArchiveFormat = "tar+gzip"
	// SkillsRepoMaxSkills bounds skills.list. Past it the listing fails closed
	// instead of returning a partial catalog.
	SkillsRepoMaxSkills = 500
	// SkillsRepoMaxFiles, SkillsRepoMaxFileBytes and SkillsRepoMaxArchiveBytes
	// bound skills.fetch. The archive bound applies to the uncompressed file
	// contents; past any bound the fetch fails closed.
	SkillsRepoMaxFiles        = 500
	SkillsRepoMaxFileBytes    = 2 << 20
	SkillsRepoMaxArchiveBytes = 8 << 20
	// SkillsRepoMaxSkillFileBytes bounds the SKILL.md read for its
	// frontmatter.
	SkillsRepoMaxSkillFileBytes = 256 << 10

	skillsRepoDefaultRef          = "main"
	skillsRepoDefaultSubdirectory = "skills/"
	skillsRepoGitHubHost          = "github.com"
)

var (
	skillName           = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)
	skillVersion        = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.+_-]{0,63}$`)
	skillHarness        = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
	gitHubOwner         = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{0,38}$`)
	gitHubRepo          = regexp.MustCompile(`^\.?[A-Za-z0-9_-][A-Za-z0-9._-]{0,99}$`)
	gitRef              = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]{0,254}$`)
	subdirectorySegment = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9._-]{0,99}$`)
)

const (
	skillsRepoMaxHarnesses      = 16
	skillsRepoMaxDescription    = 1024
	skillsRepoMaxSubdirSegments = 8
)

// SkillsRepoSkill is one skill in a listing. Only these fields leave the
// repository: the rest of SKILL.md, and every other file, is returned only by
// skills.fetch.
type SkillsRepoSkill struct {
	Name        string   `json:"name"`
	Version     string   `json:"version,omitempty"`
	Description string   `json:"description,omitempty"`
	Harnesses   []string `json:"harnesses"`
	// Path is the skill folder's path in the repository.
	Path string `json:"path"`
}

// SkillsRepoListing is the skills.list result. Commit is the commit the ref
// resolved to, so a client can fetch every skill from the same snapshot.
type SkillsRepoListing struct {
	Ref    string            `json:"ref"`
	Commit string            `json:"commit"`
	Skills []SkillsRepoSkill `json:"skills"`
}

// SkillsRepoArchive is the skills.fetch result: one skill folder as an
// archive, inline. Digest is "sha256:<hex>" over Archive, and Size is
// len(Archive); a client verifies both before unpacking.
type SkillsRepoArchive struct {
	Skill   SkillsRepoSkill `json:"skill"`
	Commit  string          `json:"commit"`
	Format  string          `json:"format"`
	Digest  string          `json:"digest"`
	Size    int             `json:"size"`
	Archive []byte          `json:"archive"`
}

// SkillsRepoSource is where a connector's skills live.
type SkillsRepoSource struct {
	Host  string
	Owner string
	Repo  string
	Ref   string
	// Subdirectory is the skills folder, with a trailing '/'; empty is the
	// repository root.
	Subdirectory string
}

// SkillPath is the repository path of a skill folder.
func (s SkillsRepoSource) SkillPath(name string) string { return s.Subdirectory + name }

// SkillsRepoSourceFromConfig reads a connector's source from its Config:
// repo_url (required), ref (default main), and subdirectory (default
// skills/, "/" for the repository root). It re-checks every value, so a
// backend never sees a traversal or an option-like ref.
func SkillsRepoSourceFromConfig(config map[string]any) (SkillsRepoSource, error) {
	raw, _ := config["repo_url"].(string)
	source, err := parseSkillsRepoURL(raw)
	if err != nil {
		return SkillsRepoSource{}, err
	}
	source.Ref = skillsRepoDefaultRef
	if value, present := config["ref"]; present {
		ref, _ := value.(string)
		if !validGitRef(ref) {
			return SkillsRepoSource{}, errors.New("skills repository ref is invalid")
		}
		source.Ref = ref
	}
	subdirectory := skillsRepoDefaultSubdirectory
	if value, present := config["subdirectory"]; present {
		subdirectory, _ = value.(string)
	}
	if source.Subdirectory, err = normalizeSubdirectory(subdirectory); err != nil {
		return SkillsRepoSource{}, err
	}
	return source, nil
}

func parseSkillsRepoURL(raw string) (SkillsRepoSource, error) {
	invalid := errors.New("skills repository URL must be https://github.com/<owner>/<repo>")
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host != skillsRepoGitHubHost || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" {
		return SkillsRepoSource{}, invalid
	}
	parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	if len(parts) != 2 {
		return SkillsRepoSource{}, invalid
	}
	owner, repo := parts[0], strings.TrimSuffix(parts[1], ".git")
	if !gitHubOwner.MatchString(owner) || !gitHubRepo.MatchString(repo) {
		return SkillsRepoSource{}, invalid
	}
	return SkillsRepoSource{Host: skillsRepoGitHubHost, Owner: owner, Repo: repo}, nil
}

// validGitRef accepts branch, tag and commit names and refuses anything git
// itself refuses that could change how the ref is read: '..', empty or
// dot-led segments, a trailing '/' or '.', '.lock' components, and '@{'.
func validGitRef(ref string) bool {
	if !gitRef.MatchString(ref) || strings.Contains(ref, "..") || strings.Contains(ref, "@{") || strings.HasSuffix(ref, ".") {
		return false
	}
	for _, segment := range strings.Split(ref, "/") {
		if segment == "" || strings.HasPrefix(segment, ".") || strings.HasSuffix(segment, ".lock") {
			return false
		}
	}
	return true
}

func normalizeSubdirectory(subdirectory string) (string, error) {
	invalid := errors.New("skills repository subdirectory is invalid")
	if subdirectory == "/" {
		return "", nil
	}
	trimmed := strings.TrimSuffix(subdirectory, "/")
	segments := strings.Split(trimmed, "/")
	if trimmed == "" || len(segments) > skillsRepoMaxSubdirSegments {
		return "", invalid
	}
	for _, segment := range segments {
		if !subdirectorySegment.MatchString(segment) {
			return "", invalid
		}
	}
	return trimmed + "/", nil
}

// ValidateSkillName checks that name is a single, traversal-free path segment.
func ValidateSkillName(name string) error {
	if !skillName.MatchString(name) {
		return errors.New("skill name must be 1-64 letters, digits, '-' or '_'")
	}
	return nil
}

var skillsRepoResources = map[string][]string{
	"skills.list":  {"skills"},
	"skills.fetch": {"skills", "*"},
}

// ValidateSkillsRepoResource checks that resource is the resource shape of a
// skills_repo capability.
func ValidateSkillsRepoResource(capability, resource string) error {
	return validateResource("skills_repo", skillsRepoResources, capability, resource, skillName.MatchString)
}

type SkillsListPayload struct{}

func (SkillsListPayload) Capability() string { return "skills.list" }
func (SkillsListPayload) Validate() error    { return nil }

// SkillsFetchPayload fetches the skill named by the resource skills/<name>.
type SkillsFetchPayload struct{}

func (SkillsFetchPayload) Capability() string { return "skills.fetch" }
func (SkillsFetchPayload) Validate() error    { return nil }

// SkillsRepoBackend is the seam for skills repository reads. FetchSkill
// returns the archive with Skill, Commit, Format and Archive set; the client
// sets Digest and Size.
type SkillsRepoBackend interface {
	ListSkills(ctx context.Context, source SkillsRepoSource) (SkillsRepoListing, error)
	FetchSkill(ctx context.Context, source SkillsRepoSource, name string) (SkillsRepoArchive, error)
}

// SkillsRepoClient serves skills.list and skills.fetch for a skills_repo
// connector, reading its source from the connector's Config.
type SkillsRepoClient struct {
	backend SkillsRepoBackend
}

func NewSkillsRepo(backend SkillsRepoBackend) *SkillsRepoClient {
	return &SkillsRepoClient{backend: backend}
}

func (c *SkillsRepoClient) Provider() contract.Provider { return contract.ProviderSkillsRepo }

func (c *SkillsRepoClient) Invoke(ctx context.Context, meta contract.Metadata, inv contract.Invocation, payload Payload) (Result, error) {
	if err := checkProvider(meta, c.Provider()); err != nil {
		return Result{}, err
	}
	if err := Guard(meta, inv); err != nil {
		return Result{}, err
	}
	if err := matchCapability(inv, payload); err != nil {
		return Result{}, err
	}
	if err := ValidateSkillsRepoResource(inv.Capability, inv.Resource); err != nil {
		return Result{}, err
	}
	source, err := SkillsRepoSourceFromConfig(meta.Config)
	if err != nil {
		return Result{}, err
	}
	ctx = withInvocationScope(ctx, meta, inv)
	switch p := payload.(type) {
	case SkillsListPayload:
		listing, err := c.backend.ListSkills(ctx, source)
		if err != nil {
			return Result{}, err
		}
		if len(listing.Skills) > SkillsRepoMaxSkills {
			return Result{}, fmt.Errorf("skills repository has more than %d skills", SkillsRepoMaxSkills)
		}
		skills := make([]SkillsRepoSkill, 0, len(listing.Skills))
		for _, skill := range listing.Skills {
			skills = append(skills, allowSkillFields(skill))
		}
		return Result{Capability: inv.Capability, Data: SkillsRepoListing{Ref: source.Ref, Commit: listing.Commit, Skills: skills}}, nil
	case SkillsFetchPayload:
		name := strings.TrimPrefix(inv.Resource, "skills/")
		archive, err := c.backend.FetchSkill(ctx, source, name)
		if err != nil {
			return Result{}, err
		}
		if archive.Skill.Name != name || archive.Format != SkillsRepoArchiveFormat || len(archive.Archive) == 0 {
			return Result{}, errors.New("skills repository returned an invalid archive")
		}
		archive.Skill = allowSkillFields(archive.Skill)
		archive.Digest, archive.Size = sha256Digest(archive.Archive), len(archive.Archive)
		return Result{Capability: inv.Capability, Data: archive}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", p)
	}
}

// allowSkillFields copies only the allowlisted listing fields, so a backend
// cannot widen what a listing returns.
func allowSkillFields(skill SkillsRepoSkill) SkillsRepoSkill {
	harnesses := append([]string{}, skill.Harnesses...)
	return SkillsRepoSkill{Name: skill.Name, Version: skill.Version, Description: skill.Description, Harnesses: harnesses, Path: skill.Path}
}

// MemorySkillsRepo is an in-memory skills repository backend for tests.
// Archives are keyed by skill name.
type MemorySkillsRepo struct {
	Commit   string
	Skills   []SkillsRepoSkill
	Archives map[string][]byte
}

func (m MemorySkillsRepo) ListSkills(context.Context, SkillsRepoSource) (SkillsRepoListing, error) {
	return SkillsRepoListing{Commit: m.Commit, Skills: m.Skills}, nil
}

func (m MemorySkillsRepo) FetchSkill(_ context.Context, _ SkillsRepoSource, name string) (SkillsRepoArchive, error) {
	archive, ok := m.Archives[name]
	if !ok {
		return SkillsRepoArchive{}, errors.New("skill not found")
	}
	for _, skill := range m.Skills {
		if skill.Name == name {
			return SkillsRepoArchive{Skill: skill, Commit: m.Commit, Format: SkillsRepoArchiveFormat, Archive: archive}, nil
		}
	}
	return SkillsRepoArchive{}, errors.New("skill not found")
}

// skillFrontmatter is the part of a SKILL.md frontmatter a listing reads.
// version and harnesses may be top-level keys or under metadata:.
type skillFrontmatter struct {
	Name        string
	Description string
	Version     string
	Harnesses   []string
}

// parseSkillFrontmatter reads the YAML frontmatter at the top of a SKILL.md.
// The module is stdlib-only, so this reads the subset skills use: scalar
// key: value pairs (optionally quoted), | and > block scalars, inline [a, b]
// lists, block "- a" lists, and one level of nesting under metadata:. Other
// keys are ignored.
func parseSkillFrontmatter(content []byte) (skillFrontmatter, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	scanner.Buffer(make([]byte, 0, 4096), SkillsRepoMaxSkillFileBytes)
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return skillFrontmatter{}, errors.New("SKILL.md has no frontmatter")
	}
	values := map[string]string{}
	lists := map[string][]string{}
	var section, listKey, blockKey, blockStyle string
	var block []string
	endBlock := func() {
		if blockKey != "" {
			separator := " "
			if strings.HasPrefix(blockStyle, "|") {
				separator = "\n"
			}
			values[blockKey] = strings.Join(block, separator)
		}
		blockKey, block = "", nil
	}
	closed := false
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " \t\r")
		if line == "---" {
			closed = true
			break
		}
		trimmed := strings.TrimSpace(line)
		indented := line != strings.TrimLeft(line, " \t")
		if blockKey != "" && (indented || trimmed == "") {
			if trimmed != "" {
				block = append(block, trimmed)
			}
			continue
		}
		endBlock()
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !indented {
			section, listKey = "", ""
		}
		if strings.HasPrefix(trimmed, "- ") && listKey != "" {
			lists[listKey] = append(lists[listKey], unquote(strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))))
			continue
		}
		key, value, ok := strings.Cut(trimmed, ":")
		if !ok || (indented && section == "") {
			continue
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if section != "" {
			key = section + "." + key
		}
		switch {
		case value == "" && !indented:
			section, listKey = key, key
		case value == "":
			listKey = key
		case strings.Trim(value, "|>+-") == "":
			blockKey, blockStyle = key, value
		case strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]"):
			for _, item := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"), ",") {
				if item = unquote(strings.TrimSpace(item)); item != "" {
					lists[key] = append(lists[key], item)
				}
			}
		default:
			values[key] = unquote(value)
		}
	}
	endBlock()
	if err := scanner.Err(); err != nil || !closed {
		return skillFrontmatter{}, errors.New("SKILL.md frontmatter is invalid")
	}
	fm := skillFrontmatter{Name: values["name"], Description: values["description"], Version: values["version"], Harnesses: lists["harnesses"]}
	if fm.Version == "" {
		fm.Version = values["metadata.version"]
	}
	if len(fm.Harnesses) == 0 {
		fm.Harnesses = lists["metadata.harnesses"]
	}
	return fm, nil
}

func unquote(value string) string {
	if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
		return value[1 : len(value)-1]
	}
	return value
}

// skillFromFrontmatter builds a listing entry from a skill folder's SKILL.md.
// The frontmatter name must equal the folder name. A version or harness that
// is not a plain token is dropped rather than returned, and the description
// is truncated to the listing bound.
func skillFromFrontmatter(name, path string, content []byte) (SkillsRepoSkill, error) {
	fm, err := parseSkillFrontmatter(content)
	if err != nil {
		return SkillsRepoSkill{}, err
	}
	if fm.Name != name {
		return SkillsRepoSkill{}, errors.New("SKILL.md name does not match its folder")
	}
	skill := SkillsRepoSkill{Name: name, Description: fm.Description, Harnesses: []string{}, Path: path}
	if skillVersion.MatchString(fm.Version) {
		skill.Version = fm.Version
	}
	if runes := []rune(skill.Description); len(runes) > skillsRepoMaxDescription {
		skill.Description = string(runes[:skillsRepoMaxDescription])
	}
	seen := map[string]bool{}
	for _, harness := range fm.Harnesses {
		if skillHarness.MatchString(harness) && !seen[harness] && len(skill.Harnesses) < skillsRepoMaxHarnesses {
			seen[harness] = true
			skill.Harnesses = append(skill.Harnesses, harness)
		}
	}
	return skill, nil
}
