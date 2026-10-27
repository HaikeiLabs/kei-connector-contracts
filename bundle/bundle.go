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

// Package bundle defines the wire types, the decoding and the validation of
// Kei policy bundles: the document the policy catalog issues to an
// installation, and the decisions every service derives from it.
//
// The package is the shared implementation of that contract. A caller that
// only reads a bundle (the CLI) and a caller that enforces one (the tenant
// runtime) therefore use the same types and the same checks, and cannot drift
// apart.
//
// The package is stdlib-only and carries no dependency of its own beyond the
// other packages of this module. Three checks are supplied by the caller
// instead, because they need a JSON Schema implementation this module does not
// carry:
//
//   - SetSchemaValidator and SetV2SchemaValidator install the validators for
//     the bundle's v1 and v2 JSON Schema documents.
//   - SetPlanValidator installs the validator for a v2 tool's operation plan,
//     which also compiles the plan's args_schema and context_schema.
//
// A validator that is not installed is not applied: the bundle is then checked
// by this package's own structural validation only. A service that enforces
// bundles installs all three.
package bundle

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/harnessmatch"
)

var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

// ContentType is the media type for a v1 policy bundle: the legacy model the
// catalog still folds its v2 policy set into, and the one the runtime still
// accepts. Kept for callers that check the v1 constant.
const ContentType = "application/vnd.kei.policy-bundle.v1+json"

// PlanValidator validates the operation plan of one v2 connector tool: the
// plan document and the tool registration it has to match. The tool
// registration matters because a plan is only valid against the capabilities
// and resource types its own tool declares.
//
// A plan is a declarative document whose args_schema and context_schema are
// JSON Schema documents; validating them needs a JSON Schema implementation,
// which this dependency-free module does not carry. The enforcing service
// supplies one and installs it with SetPlanValidator.
type PlanValidator func(tool V2Tool) error

// ContentTypeV2 is the media type for the current v2 policy bundle.
const ContentTypeV2 = "application/vnd.kei.policy-bundle.v2+json"

// Issuer is the expected issuer claim in every valid bundle.
const Issuer = "kei-policy-catalog"

// MaxIssuedAtClockSkew is how far a bundle's issued_at may be ahead of the
// runtime clock. The catalog stamps issued_at with its own clock, which can
// lead a runtime host by a few hundred milliseconds.
const MaxIssuedAtClockSkew = 5 * time.Minute

var (
	ErrBundleExpired        = errors.New("policy bundle expired")
	ErrBundleUnsupported    = errors.New("policy bundle unsupported")
	ErrSubjectSnapshotStale = errors.New("policy bundle subject snapshot is stale")
)

// Audience identifies the installation a bundle is scoped to.
type Audience struct {
	InstallationID string `json:"installation_id"`
	OrgID          string `json:"org_id"`
	WorkspaceID    string `json:"workspace_id"`
}

// Identity is one known user in a bundle with subject_resolution.mode=bundle.
type Identity struct {
	Source       string    `json:"source,omitempty"`
	ExternalID   string    `json:"external_id,omitempty"`
	UserID       string    `json:"user_id"`
	Status       string    `json:"status"`
	Groups       []string  `json:"groups,omitempty"`
	OrgRoles     []string  `json:"org_roles,omitempty"`
	SubjectToken string    `json:"subject_token,omitempty"`
	ExpireTime   time.Time `json:"expire_time,omitempty"`
}

// Harness is one registered coding-agent harness (ADR-029 §2).
type Harness struct {
	AgentID string `json:"agent_id"`
	Kind    string `json:"kind"`
}

// HarnessKind is one entry of the bundle's harness kind registry (CA-5): how a
// kind is named in hooks, the prefix of its tool names, the tool that carries a
// shell command line, and the native tools a policy can target. Gateable is
// false for a kind that can only be observed, never gated.
type HarnessKind struct {
	Kind       string `json:"kind"`
	ShortName  string `json:"short_name,omitempty"`
	ToolPrefix string `json:"tool_prefix,omitempty"`
	ShellTool  string `json:"shell_tool,omitempty"`
	Gateable   bool   `json:"gateable"`
	// DecidesVia says how the kind's calls are decided: DecidesViaPDP for the
	// kei-proxy local PDP, DecidesViaNative for a rendered native allowlist.
	// Empty keeps the built-in meaning (custom and pi via the PDP).
	DecidesVia  string       `json:"decides_via,omitempty"`
	NativeTools []NativeTool `json:"native_tools,omitempty"`
}

// Values of HarnessKind.DecidesVia.
const (
	DecidesViaPDP    = "pdp"
	DecidesViaNative = "native"
)

// NativeTool maps the registry tool name a policy targets to the name the
// harness itself uses for it.
type NativeTool struct {
	ToolName   string `json:"tool_name"`
	NativeName string `json:"native_name"`
}

// Payload is the signed-envelope body of a policy bundle.
type Payload struct {
	Schema         string    `json:"schema"`
	BundleID       string    `json:"bundle_id"`
	BundleVersion  int64     `json:"bundle_version"`
	PolicyRevision int64     `json:"policy_revision"`
	Issuer         string    `json:"issuer"`
	Audience       Audience  `json:"audience"`
	IssuedAt       time.Time `json:"issued_at"`
	NotAfter       time.Time `json:"not_after"`
	Refresh        struct {
		PollIntervalSeconds int `json:"poll_interval_seconds"`
		StaleAfterSeconds   int `json:"stale_after_seconds"`
	} `json:"refresh"`
	SubjectResolution struct {
		Mode            string `json:"mode"`
		MaxCacheSeconds int    `json:"max_cache_seconds"`
	} `json:"subject_resolution"`
	PolicySet              json.RawMessage `json:"policy_set"`
	AgentConnectorBindings json.RawMessage `json:"agent_connector_bindings"`
	Identities             []Identity      `json:"identities,omitempty"`
	Harnesses              []Harness       `json:"harnesses,omitempty"`
	// HarnessKinds is the harness kind registry. When absent, only the
	// built-in kinds (harnessmatch.Kinds) are valid.
	HarnessKinds          []HarnessKind `json:"harness_kinds,omitempty"`
	HarnessMatchSemantics string        `json:"harness_match_semantics,omitempty"`
	// Subject is the v2-only installation subject (HAI-433).
	Subject *InstallationSubject `json:"subject,omitempty"`
	// AgentModelProfiles is the v2-only effective model profile of each agent
	// assigned to the installation (HAI-459). In v2 it is the model-call
	// permission; see EvaluateV2Model.
	AgentModelProfiles []AgentModelProfile `json:"agent_model_profiles,omitempty"`
}

// AgentModelProfile is one agent's effective model profile in a v2 bundle
// (ADR-035 §6). Status is the profile's catalog status; only active permits.
type AgentModelProfile struct {
	AgentID   string `json:"agent_id"`
	ProfileID string `json:"profile_id"`
	Status    string `json:"status"`
}

// InstallationSubject is the installation's principal as the catalog resolved
// it (HAI-433): the installation creator, their email, and their effective
// group names. Only kei.policy-bundle/v2 carries it; harness renderers apply
// user:, email: and group: harness_policies sources against it.
type InstallationSubject struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email,omitempty"`
	Groups []string `json:"groups"`
}

// Verified wraps a successfully validated bundle together with its raw bytes
// and digest.
type Verified struct {
	Payload      Payload
	PayloadBytes []byte
	Digest       string
}

// Validate checks an unsigned catalog bundle retrieved over the authenticated
// HTTPS runtime-control-plane connection. The digest binds the exact received
// bytes to the catalog ETag and persisted state; it is not a signature.
func Validate(bundleBytes []byte, installationID, orgID, workspaceID string, highestVersion int64, now time.Time) (*Verified, error) {
	if len(bundleBytes) == 0 {
		return nil, errors.New("empty policy bundle")
	}
	if err := RejectDuplicateJSON(bundleBytes); err != nil {
		return nil, err
	}
	var identity struct {
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal(bundleBytes, &identity); err != nil {
		return nil, fmt.Errorf("decode policy bundle identity: %w", err)
	}
	if identity.Schema == "kei.policy-bundle/v2" {
		return ValidateV2(bundleBytes, installationID, orgID, workspaceID, highestVersion, now)
	}
	if identity.Schema != "kei.policy-bundle/v1" {
		return nil, fmt.Errorf("%w: schema", ErrBundleUnsupported)
	}
	if err := validatePayloadSchema(bundleBytes); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "additional properties") {
			return nil, fmt.Errorf("%w: closed schema rejected unknown field", ErrBundleUnsupported)
		}
		return nil, fmt.Errorf("invalid policy bundle schema: %w", err)
	}
	var payload Payload
	if err := strictJSON(bundleBytes, &payload); err != nil {
		return nil, fmt.Errorf("decode policy bundle: %w", err)
	}
	if payload.Issuer != Issuer || payload.BundleID == "" || payload.BundleVersion < 1 || payload.PolicyRevision < 1 {
		return nil, errors.New("invalid policy bundle identity")
	}
	if payload.Audience.InstallationID != installationID || payload.Audience.OrgID != orgID || payload.Audience.WorkspaceID != workspaceID {
		return nil, errors.New("policy bundle audience mismatch")
	}
	if payload.BundleVersion <= highestVersion {
		return nil, errors.New("policy bundle version rollback or replay")
	}
	if payload.IssuedAt.After(now.Add(MaxIssuedAtClockSkew)) {
		return nil, errors.New("policy bundle issued too far in the future")
	}
	if !now.Before(payload.NotAfter) {
		return nil, ErrBundleExpired
	}
	if payload.NotAfter.After(payload.IssuedAt.Add(24*time.Hour)) || !payload.NotAfter.After(payload.IssuedAt) {
		return nil, errors.New("invalid policy bundle validity window")
	}
	switch payload.SubjectResolution.Mode {
	case "online":
	case "bundle":
	default:
		return nil, fmt.Errorf("%w: subject_resolution.mode", ErrBundleUnsupported)
	}
	if err := validatePolicyDialect(payload.PolicySet); err != nil {
		return nil, err
	}
	if err := validateHarnesses(payload); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBundleUnsupported, err)
	}
	digest := sha256.Sum256(bundleBytes)
	return &Verified{Payload: payload, PayloadBytes: append([]byte(nil), bundleBytes...), Digest: fmt.Sprintf("sha256:%x", digest)}, nil
}

var (
	schemaValidator   func([]byte) error
	schemaValidatorV2 func([]byte) error
	planValidator     PlanValidator
)

// SetSchemaValidator installs the validator for the v1 bundle's JSON Schema
// document. It is the catalog's schema, compiled by the caller. A nil
// validator leaves the v1 schema unchecked: Validate then applies this
// package's own structural checks only.
func SetSchemaValidator(v func([]byte) error) { schemaValidator = v }

// SetV2SchemaValidator installs the validator for the v2 bundle's JSON Schema
// document. It is the catalog's schema, compiled by the caller. A nil
// validator leaves the v2 schema unchecked: ValidateV2 then applies this
// package's own structural checks only.
func SetV2SchemaValidator(v func([]byte) error) { schemaValidatorV2 = v }

// SetPlanValidator installs the validator for a v2 connector tool's
// operation plan. The enforcing service sets it to the validator that
// enforces the plan profile; until one is set, plans are not validated.
func SetPlanValidator(v PlanValidator) { planValidator = v }

// planValid returns nil when a tool's operation plan is valid under the
// installed plan validator. With no validator installed the plan is not
// checked at all.
func planValid(tool V2Tool) error {
	if planValidator == nil {
		return nil
	}
	return planValidator(tool)
}

// ValidateV2 validates the ADR-035 v2 envelope. It is separate from Validate
// so that a caller can accept the current bundle format without changing the
// behaviour of the v1 path.
func ValidateV2(bundleBytes []byte, installationID, orgID, workspaceID string, highestVersion int64, now time.Time) (*Verified, error) {
	if len(bundleBytes) == 0 {
		return nil, errors.New("empty policy bundle")
	}
	if err := RejectDuplicateJSON(bundleBytes); err != nil {
		return nil, err
	}
	if schemaValidatorV2 == nil {
		return nil, errors.New("policy bundle v2 schema validator is unavailable")
	}
	if err := schemaValidatorV2(bundleBytes); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "additional properties") {
			return nil, fmt.Errorf("%w: closed v2 schema rejected unknown field", ErrBundleUnsupported)
		}
		return nil, fmt.Errorf("invalid policy bundle v2 schema: %w", err)
	}
	var payload Payload
	if err := strictJSON(bundleBytes, &payload); err != nil {
		return nil, fmt.Errorf("decode policy bundle v2: %w", err)
	}
	if payload.Schema != "kei.policy-bundle/v2" || payload.Issuer != Issuer || payload.BundleID == "" || payload.BundleVersion < 1 || payload.PolicyRevision < 1 {
		return nil, errors.New("invalid policy bundle v2 identity")
	}
	if payload.Audience.InstallationID != installationID || payload.Audience.OrgID != orgID || payload.Audience.WorkspaceID != workspaceID {
		return nil, errors.New("policy bundle v2 audience mismatch")
	}
	if payload.BundleVersion <= highestVersion {
		return nil, errors.New("policy bundle v2 version rollback or replay")
	}
	if payload.IssuedAt.After(now.Add(MaxIssuedAtClockSkew)) {
		return nil, errors.New("policy bundle v2 issued too far in the future")
	}
	if !now.Before(payload.NotAfter) {
		return nil, ErrBundleExpired
	}
	if payload.NotAfter.After(payload.IssuedAt.Add(24*time.Hour)) || !payload.NotAfter.After(payload.IssuedAt) {
		return nil, errors.New("invalid policy bundle v2 validity window")
	}
	if payload.SubjectResolution.Mode != "bundle" {
		return nil, fmt.Errorf("%w: subject_resolution.mode", ErrBundleUnsupported)
	}
	if payload.Identities == nil {
		return nil, fmt.Errorf("%w: identities", ErrBundleUnsupported)
	}
	if err := validatePolicySetV2(payload.PolicySet); err != nil {
		return nil, err
	}
	if err := validateHarnesses(payload); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBundleUnsupported, err)
	}
	digest := sha256.Sum256(bundleBytes)
	return &Verified{Payload: payload, PayloadBytes: append([]byte(nil), bundleBytes...), Digest: fmt.Sprintf("sha256:%x", digest)}, nil
}

func validatePolicySetV2(raw json.RawMessage) error {
	var v struct {
		Schema         string `json:"schema"`
		MatchSemantics string `json:"match_semantics"`
		Evaluation     struct {
			Checks  string `json:"checks"`
			Default string `json:"default"`
		} `json:"evaluation"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	if v.Schema != "kei.policy-set/v2" {
		return fmt.Errorf("%w: policy_set.schema", ErrBundleUnsupported)
	}
	if v.MatchSemantics != "kei.match/v2" {
		return fmt.Errorf("%w: policy_set.match_semantics", ErrBundleUnsupported)
	}
	if v.Evaluation.Checks != "per_capability_first_match_and" {
		return fmt.Errorf("%w: policy_set.evaluation.checks", ErrBundleUnsupported)
	}
	if v.Evaluation.Default != "deny" {
		return fmt.Errorf("%w: policy_set.evaluation.default", ErrBundleUnsupported)
	}
	var set v2PolicySet
	if err := strictJSON(raw, &set); err != nil {
		return fmt.Errorf("%w: policy_set", ErrBundleUnsupported)
	}
	for _, tool := range set.Tools {
		if tool.Route.ConnectorBinding != nil && len(tool.Plan) > 0 {
			if err := planValid(tool); err != nil {
				return fmt.Errorf("%w: tool %s: %v", ErrBundleUnsupported, tool.Name, err)
			}
		}
	}
	return nil
}

func validatePayloadSchema(raw []byte) error {
	if err := RejectDuplicateJSON(raw); err != nil {
		return err
	}
	if schemaValidator == nil {
		return errors.New("policy bundle schema validator is unavailable")
	}
	return schemaValidator(raw)
}

func validatePolicyDialect(raw json.RawMessage) error {
	var v struct {
		MatchSemantics string `json:"match_semantics"`
		Evaluation     struct {
			Order   string `json:"order"`
			Default string `json:"default"`
		} `json:"evaluation"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	if v.MatchSemantics != "kei.match/v1" {
		return fmt.Errorf("%w: policy_set.match_semantics", ErrBundleUnsupported)
	}
	if v.Evaluation.Order != "array_order_first_match" {
		return fmt.Errorf("%w: policy_set.evaluation.order", ErrBundleUnsupported)
	}
	if v.Evaluation.Default != "deny" {
		return fmt.Errorf("%w: policy_set.evaluation.default", ErrBundleUnsupported)
	}
	return nil
}

// KnownHarnessKind reports whether kind is a harness kind of this bundle: one
// the bundle's registry lists, or a built-in kind when it carries no registry.
func (p Payload) KnownHarnessKind(kind string) bool {
	if len(p.HarnessKinds) == 0 {
		return harnessmatch.ValidKind(kind)
	}
	for _, k := range p.HarnessKinds {
		if k.Kind == kind {
			return true
		}
	}
	return false
}

// HarnessKindDecidesViaPDP reports whether kind decides each call through the
// local PDP: the bundle's registry entry says so, and a kind it does not
// describe keeps the built-in meaning (custom and pi).
func (p Payload) HarnessKindDecidesViaPDP(kind string) bool {
	for _, k := range p.HarnessKinds {
		if k.Kind == kind && k.DecidesVia != "" {
			return k.DecidesVia == DecidesViaPDP
		}
	}
	return kind == "custom" || kind == "pi"
}

// legacyShellKind is the shell tool of the Pi kind: the released adapter name,
// honoured while a bundle carries no registry.
var legacyShellKind = HarnessKind{Kind: "pi", ShellTool: "pi.bash", DecidesVia: DecidesViaPDP, Gateable: true}

// HarnessKindForShellTool returns the kind whose shell command line the
// local PDP decides when tool is its shell tool. A bundle with no registry
// knows only the legacy pi.bash; with a registry only its own entries count.
func (p Payload) HarnessKindForShellTool(tool string) (HarnessKind, bool) {
	if tool == "" {
		return HarnessKind{}, false
	}
	if len(p.HarnessKinds) == 0 {
		return legacyShellKind, tool == legacyShellKind.ShellTool
	}
	for _, k := range p.HarnessKinds {
		if k.ShellTool == tool && p.HarnessKindDecidesViaPDP(k.Kind) {
			return k, true
		}
	}
	return HarnessKind{}, false
}

var harnessKindPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)

func validateHarnessKinds(kinds []HarnessKind) error {
	seen := map[string]bool{}
	for i, k := range kinds {
		if !harnessKindPattern.MatchString(k.Kind) {
			return fmt.Errorf("harness_kinds[%d].kind %q is not a valid harness kind", i, k.Kind)
		}
		if seen[k.Kind] {
			return fmt.Errorf("harness_kinds[%d].kind %q is listed twice", i, k.Kind)
		}
		seen[k.Kind] = true
		switch k.DecidesVia {
		case "", DecidesViaPDP, DecidesViaNative:
		default:
			return fmt.Errorf("harness_kinds[%d].decides_via %q must be %q or %q", i, k.DecidesVia, DecidesViaPDP, DecidesViaNative)
		}
		if k.ShellTool != "" && !strings.HasPrefix(k.ShellTool, k.ToolPrefix) {
			return fmt.Errorf("harness_kinds[%d].shell_tool %q must start with tool_prefix %q", i, k.ShellTool, k.ToolPrefix)
		}
		for j, n := range k.NativeTools {
			if k.ToolPrefix == "" || !strings.HasPrefix(n.ToolName, k.ToolPrefix) || n.NativeName == "" {
				return fmt.Errorf("harness_kinds[%d].native_tools[%d] needs a tool_name under tool_prefix %q and a native_name", i, j, k.ToolPrefix)
			}
		}
	}
	return nil
}

func validateHarnesses(payload Payload) error {
	harnesses, semantics := payload.Harnesses, payload.HarnessMatchSemantics
	if err := validateHarnessKinds(payload.HarnessKinds); err != nil {
		return err
	}
	if len(harnesses) == 0 && semantics == "" {
		return nil
	}
	if semantics != "kei.harness-match/v1" {
		return fmt.Errorf("harness_match_semantics must be %q", "kei.harness-match/v1")
	}
	for i, h := range harnesses {
		if !payload.KnownHarnessKind(h.Kind) {
			return fmt.Errorf("harnesses[%d].kind %q is not a valid harness kind", i, h.Kind)
		}
		if !idPattern.MatchString(h.AgentID) {
			return fmt.Errorf("harnesses[%d].agent_id %q is not a valid identifier", i, h.AgentID)
		}
		if payload.KnownHarnessKind(h.AgentID) {
			return fmt.Errorf("harnesses[%d].agent_id %q must not equal a harness kind value", i, h.AgentID)
		}
	}
	return nil
}

func strictJSON(raw []byte, dst any) error {
	if err := RejectDuplicateJSON(raw); err != nil {
		return err
	}
	d := json.NewDecoder(strings.NewReader(string(raw)))
	d.DisallowUnknownFields()
	if err := d.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("trailing JSON value")
		}
		return err
	}
	return nil
}

// RejectDuplicateJSON prevents ambiguous last-key-wins parsing of bundle bytes.
func RejectDuplicateJSON(raw []byte) error {
	d := json.NewDecoder(strings.NewReader(string(raw)))
	d.UseNumber()
	var value func() error
	value = func() error {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		delim, ok := tok.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				keyTok, err := d.Token()
				if err != nil {
					return err
				}
				key, ok := keyTok.(string)
				if !ok {
					return errors.New("invalid JSON object key")
				}
				if seen[key] {
					return fmt.Errorf("duplicate JSON object key %q", key)
				}
				seen[key] = true
				if err := value(); err != nil {
					return err
				}
			}
			end, err := d.Token()
			if err != nil || end != json.Delim('}') {
				return errors.New("malformed JSON object")
			}
		case '[':
			for d.More() {
				if err := value(); err != nil {
					return err
				}
			}
			end, err := d.Token()
			if err != nil || end != json.Delim(']') {
				return errors.New("malformed JSON array")
			}
		default:
			return errors.New("unexpected JSON delimiter")
		}
		return nil
	}
	if err := value(); err != nil {
		return err
	}
	if _, err := d.Token(); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("trailing JSON value")
		}
		return err
	}
	return nil
}
