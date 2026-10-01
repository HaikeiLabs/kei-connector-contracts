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
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

func grafanaPayloads() []Payload {
	return []Payload{
		GrafanaFolderListPayload{}, GrafanaFolderReadPayload{}, GrafanaDashboardSearchPayload{},
		GrafanaDashboardReadPayload{}, GrafanaDatasourceListPayload{}, GrafanaDatasourceReadPayload{},
		GrafanaAlertRuleListPayload{}, GrafanaAnnotationListPayload{}, GrafanaDatasourceQueryPayload{},
	}
}

func TestGrafanaPayloadsCoverTheCapabilities(t *testing.T) {
	assertPayloadsCover(t, contract.ProviderGrafana, grafanaPayloads())
}

func TestGrafanaResourceGrammar(t *testing.T) {
	valid := map[string]string{
		"folder.list":      "folders",
		"folder.read":      "folders/fd-ops",
		"dashboard.search": "search",
		"dashboard.read":   "dashboards/cIBgcSjkk",
		"datasource.list":  "datasources",
		"datasource.read":  "datasources/P8E80F9AEF21F6940",
		"alert_rule.list":  "alert_rules",
		"annotation.list":  "annotations",
		"datasource.query": "datasources/P8E80F9AEF21F6940/query",
	}
	for capability, resource := range valid {
		if err := ValidateGrafanaResource(capability, resource); err != nil {
			t.Errorf("%s %q rejected: %v", capability, resource, err)
		}
	}
	for _, tc := range []struct{ capability, resource string }{
		{"folder.list", "folders/fd-ops"},
		{"folder.read", "folders"},
		{"folder.read", "folders/has space"},
		{"folder.read", "folders/" + strings.Repeat("a", 41)},
		{"dashboard.search", "search/cpu"},
		{"dashboard.read", "dashboards/uid/cIBgcSjkk"},
		{"datasource.read", "datasources/P8E80F9AEF21F6940/query"},
		{"datasource.query", "datasources/P8E80F9AEF21F6940"},
		{"datasource.query", "datasources/../query"},
		{"alert_rule.list", "alert_rules/r-1"},
		{"annotation.list", "annotations/1"},
		{"dashboard.create", "dashboards/cIBgcSjkk"},
		{"dashboard.read", "/dashboards/cIBgcSjkk"},
	} {
		if err := ValidateGrafanaResource(tc.capability, tc.resource); err == nil {
			t.Errorf("%s %q accepted", tc.capability, tc.resource)
		}
	}
}

func TestGrafanaDashboardSearchInputIsBounded(t *testing.T) {
	for _, p := range []GrafanaDashboardSearchPayload{
		{},
		{Query: "cpu", Tags: []string{"prod"}, Limit: 100},
	} {
		if err := p.Validate(); err != nil {
			t.Errorf("%+v rejected: %v", p, err)
		}
	}
	for _, p := range []GrafanaDashboardSearchPayload{
		{Limit: -1},
		{Limit: 101},
		{Query: strings.Repeat("q", 257)},
		{Tags: make([]string, 11)},
		{Tags: []string{""}},
		{Tags: []string{strings.Repeat("t", 129)}},
	} {
		if err := p.Validate(); err == nil {
			t.Errorf("%+v accepted", p)
		}
	}
}

func TestGrafanaAnnotationListRequiresABoundedWindow(t *testing.T) {
	to := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	for _, p := range []GrafanaAnnotationListPayload{
		{From: to.Add(-time.Hour), To: to},
		{From: to.Add(-GrafanaMaxAnnotationWindow), To: to, Limit: 100, DashboardUID: "cIBgcSjkk", Tags: []string{"deploy"}},
	} {
		if err := p.Validate(); err != nil {
			t.Errorf("%+v rejected: %v", p, err)
		}
	}
	for _, p := range []GrafanaAnnotationListPayload{
		{},
		{To: to},
		{From: to},
		{From: to, To: to},
		{From: to, To: to.Add(-time.Hour)},
		{From: to.Add(-GrafanaMaxAnnotationWindow - time.Second), To: to},
		{From: to.Add(-time.Hour), To: to, Limit: 101},
		{From: to.Add(-time.Hour), To: to, DashboardUID: "bad uid"},
	} {
		if err := p.Validate(); err == nil {
			t.Errorf("%+v accepted", p)
		}
	}
}

func TestGrafanaDatasourceQueryIsBoundedInTimeAndRows(t *testing.T) {
	to := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	query := []map[string]any{{"refId": "A", "expr": "up"}}
	for _, p := range []GrafanaDatasourceQueryPayload{
		{From: to.Add(-time.Hour), To: to, Queries: query},
		{From: to.Add(-GrafanaMaxQueryWindow), To: to, MaxRows: GrafanaMaxQueryRows, Queries: query},
	} {
		if err := p.Validate(); err != nil {
			t.Errorf("%+v rejected: %v", p, err)
		}
	}
	for _, p := range []GrafanaDatasourceQueryPayload{
		{Queries: query},
		{From: to.Add(-time.Hour), To: to},
		{From: to, To: to, Queries: query},
		{From: to.Add(-GrafanaMaxQueryWindow - time.Second), To: to, Queries: query},
		{From: to.Add(-time.Hour), To: to, MaxRows: -1, Queries: query},
		{From: to.Add(-time.Hour), To: to, MaxRows: GrafanaMaxQueryRows + 1, Queries: query},
		{From: to.Add(-time.Hour), To: to, Queries: make([]map[string]any, 11)},
		{From: to.Add(-time.Hour), To: to, Queries: []map[string]any{nil}},
	} {
		if err := p.Validate(); err == nil {
			t.Errorf("%+v accepted", p)
		}
	}
	if got := (GrafanaDatasourceQueryPayload{}).EffectiveMaxRows(); got != GrafanaMaxQueryRows {
		t.Errorf("default max rows = %d, want %d", got, GrafanaMaxQueryRows)
	}
}
