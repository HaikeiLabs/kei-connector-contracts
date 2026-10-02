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

// Copyright 2026 Haikei Labs
// SPDX-License-Identifier: Apache-2.0

// Package policybundlehealth defines the shared runtime policy-bundle health
// wire contract. It contains metadata only; bundle bytes and tenant content
// remain local to the runtime.
package policybundlehealth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"
)

type State string

const (
	StateCold        State = "cold"
	StateActive      State = "active"
	StateStale       State = "stale_but_valid"
	StateExpired     State = "expired"
	StateInvalid     State = "invalid"
	StateUnsupported State = "unsupported"
	StateRevoked     State = "revoked"
)

type Reason string

const (
	ReasonBundleMissing     Reason = "bundle_missing"
	ReasonRefreshFailed     Reason = "refresh_failed"
	ReasonBundleStale       Reason = "bundle_stale"
	ReasonBundleExpired     Reason = "bundle_expired"
	ReasonIntegrityRejected Reason = "integrity_rejected"
	ReasonAudienceRejected  Reason = "audience_rejected"
	ReasonSchemaUnsupported Reason = "schema_unsupported"
	ReasonRollbackRejected  Reason = "rollback_rejected"
	ReasonRuntimeRevoked    Reason = "runtime_revoked"
)

// Health is the nested heartbeat report and read projection. ReportedAt is
// catalog-generated and is omitted from heartbeat requests.
type Health struct {
	SchemaVersion  int        `json:"schema_version"`
	State          State      `json:"state"`
	BundleVersion  *int64     `json:"bundle_version"`
	PolicyRevision *int64     `json:"policy_revision"`
	BundleDigest   *string    `json:"bundle_digest"`
	CheckedAt      *time.Time `json:"checked_at"`
	AcceptedAt     *time.Time `json:"accepted_at"`
	ExpiresAt      *time.Time `json:"expires_at"`
	ReasonCode     *Reason    `json:"reason_code"`
	ReportedAt     *time.Time `json:"reported_at,omitempty"`
}

var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// Validate checks the approved v1 closed-enum, nullable-field, and
// state/reason invariants. It intentionally returns generic errors only.
func (h Health) Validate(readProjection bool) error {
	if h.SchemaVersion != 1 {
		return errors.New("unsupported policy bundle health schema")
	}
	if !knownState(h.State) {
		return errors.New("invalid policy bundle health state")
	}
	if readProjection != (h.ReportedAt != nil) {
		return errors.New("invalid policy bundle health receipt timestamp")
	}
	metadata := []*int64{h.BundleVersion, h.PolicyRevision}
	present := 0
	for _, value := range metadata {
		if value != nil {
			present++
			if *value < 0 {
				return errors.New("invalid policy bundle health version")
			}
		}
	}
	if h.BundleDigest != nil {
		present++
	}
	if h.AcceptedAt != nil {
		present++
	}
	if h.ExpiresAt != nil {
		present++
	}
	if present != 0 && present != 5 {
		return errors.New("incomplete accepted bundle metadata")
	}
	if h.BundleDigest != nil && !digestPattern.MatchString(*h.BundleDigest) {
		return errors.New("invalid policy bundle digest")
	}
	for _, timestamp := range []*time.Time{h.CheckedAt, h.AcceptedAt, h.ExpiresAt, h.ReportedAt} {
		if timestamp != nil {
			_, offset := timestamp.Zone()
			if offset != 0 || timestamp.IsZero() {
				return errors.New("invalid policy bundle health timestamp")
			}
		}
	}
	if (h.CheckedAt == nil) != (h.State == StateCold && h.ReasonCode == nil) {
		// Cold is the only state allowed to omit a completed check; non-cold
		// states always describe a completed observation.
		if h.CheckedAt == nil || h.State != StateCold {
			return errors.New("invalid policy bundle check timestamp")
		}
	}
	if h.AcceptedAt != nil {
		if !h.AcceptedAt.Before(*h.ExpiresAt) {
			return errors.New("invalid policy bundle validity interval")
		}
		if h.CheckedAt == nil || h.CheckedAt.Before(*h.AcceptedAt) {
			return errors.New("invalid policy bundle check ordering")
		}
	}
	if !validStateReason(h.State, h.ReasonCode, present == 5, h.CheckedAt != nil) {
		return errors.New("invalid policy bundle state and reason combination")
	}
	return nil
}

func knownState(s State) bool {
	switch s {
	case StateCold, StateActive, StateStale, StateExpired, StateInvalid, StateUnsupported, StateRevoked:
		return true
	}
	return false
}

func validStateReason(state State, reason *Reason, accepted, checked bool) bool {
	is := func(want Reason) bool { return reason != nil && *reason == want }
	switch state {
	case StateCold:
		return !accepted && (checked && is(ReasonBundleMissing) || !checked && reason == nil)
	case StateActive:
		return accepted && checked && reason == nil
	case StateStale:
		return accepted && checked && is(ReasonBundleStale)
	case StateExpired:
		return accepted && checked && is(ReasonBundleExpired)
	case StateInvalid:
		return checked && (is(ReasonRefreshFailed) || is(ReasonIntegrityRejected) || is(ReasonAudienceRejected) || is(ReasonRollbackRejected))
	case StateUnsupported:
		return checked && is(ReasonSchemaUnsupported)
	case StateRevoked:
		return checked && is(ReasonRuntimeRevoked)
	default:
		return false
	}
}

// Decode decodes one closed health object, rejects unknown fields and trailing
// values, and validates report-versus-projection timestamp presence.
func Decode(raw []byte, readProjection bool) (Health, error) {
	var h Health
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&h); err != nil {
		return Health{}, errors.New("invalid policy bundle health report")
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return Health{}, errors.New("invalid policy bundle health report")
	}
	if err := h.Validate(readProjection); err != nil {
		return Health{}, errors.New("invalid policy bundle health report")
	}
	return h, nil
}

// ValidateReadProjection validates a catalog-enriched health snapshot.
func ValidateReadProjection(h Health) error {
	if err := h.Validate(true); err != nil {
		return fmt.Errorf("invalid read projection: %w", err)
	}
	return nil
}
