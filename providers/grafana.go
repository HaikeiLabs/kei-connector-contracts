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
	"errors"
	"fmt"
	"regexp"
	"time"
)

// Grafana is a contract-only provider: the payloads and resource grammar
// below are what the tenant runtime executes against the instance in
// Metadata.Grafana.BaseURL with a service-account token resolved from the
// connector's opaque_ref. There is no backend seam here; the runtime
// implements the reads directly.
//
// Resources (uid = Grafana UID, 1-40 of [A-Za-z0-9_-]):
//
//	folder.list       folders
//	folder.read       folders/<uid>
//	dashboard.search  search
//	dashboard.read    dashboards/<uid>
//	datasource.list   datasources
//	datasource.read   datasources/<uid>
//	alert_rule.list   alert_rules
//	annotation.list   annotations
//	datasource.query  datasources/<uid>/query

const (
	// GrafanaMaxSearchLimit bounds dashboard.search and annotation.list pages.
	GrafanaMaxSearchLimit = 100
	// GrafanaMaxAnnotationWindow is the widest from/to range annotation.list
	// may request.
	GrafanaMaxAnnotationWindow = 31 * 24 * time.Hour
	// GrafanaMaxQueryWindow is the widest from/to range datasource.query may
	// request.
	GrafanaMaxQueryWindow = 7 * 24 * time.Hour
	// GrafanaMaxQueryRows is the most rows datasource.query returns, and the
	// default when MaxRows is 0.
	GrafanaMaxQueryRows = 1000
	// GrafanaMaxQueries is the most queries one datasource.query may carry.
	GrafanaMaxQueries = 10
)

var grafanaUID = regexp.MustCompile(`^[A-Za-z0-9_-]{1,40}$`)

var grafanaResources = map[string][]string{
	"folder.list":      {"folders"},
	"folder.read":      {"folders", "*"},
	"dashboard.search": {"search"},
	"dashboard.read":   {"dashboards", "*"},
	"datasource.list":  {"datasources"},
	"datasource.read":  {"datasources", "*"},
	"alert_rule.list":  {"alert_rules"},
	"annotation.list":  {"annotations"},
	"datasource.query": {"datasources", "*", "query"},
}

// ValidateGrafanaResource checks that resource is the resource shape of a
// Grafana capability.
func ValidateGrafanaResource(capability, resource string) error {
	return validateResource("grafana", grafanaResources, capability, resource, grafanaUID.MatchString)
}

type GrafanaFolderListPayload struct{}

func (GrafanaFolderListPayload) Capability() string { return "folder.list" }
func (GrafanaFolderListPayload) Validate() error    { return nil }

type GrafanaFolderReadPayload struct{}

func (GrafanaFolderReadPayload) Capability() string { return "folder.read" }
func (GrafanaFolderReadPayload) Validate() error    { return nil }

// GrafanaDashboardSearchPayload searches dashboards by title and tags. Limit
// 0 is the provider default.
type GrafanaDashboardSearchPayload struct {
	Query string   `json:"query,omitempty"`
	Tags  []string `json:"tags,omitempty"`
	Limit int      `json:"limit,omitempty"`
}

func (GrafanaDashboardSearchPayload) Capability() string { return "dashboard.search" }

func (p GrafanaDashboardSearchPayload) Validate() error {
	if len(p.Query) > 256 {
		return errors.New("dashboard.search query is too long")
	}
	if err := validateGrafanaTags("dashboard.search", p.Tags); err != nil {
		return err
	}
	return validateGrafanaLimit("dashboard.search", p.Limit)
}

type GrafanaDashboardReadPayload struct{}

func (GrafanaDashboardReadPayload) Capability() string { return "dashboard.read" }
func (GrafanaDashboardReadPayload) Validate() error    { return nil }

type GrafanaDatasourceListPayload struct{}

func (GrafanaDatasourceListPayload) Capability() string { return "datasource.list" }
func (GrafanaDatasourceListPayload) Validate() error    { return nil }

type GrafanaDatasourceReadPayload struct{}

func (GrafanaDatasourceReadPayload) Capability() string { return "datasource.read" }
func (GrafanaDatasourceReadPayload) Validate() error    { return nil }

type GrafanaAlertRuleListPayload struct{}

func (GrafanaAlertRuleListPayload) Capability() string { return "alert_rule.list" }
func (GrafanaAlertRuleListPayload) Validate() error    { return nil }

// GrafanaAnnotationListPayload lists annotations in a required, bounded time
// window, optionally narrowed to one dashboard or to tags.
type GrafanaAnnotationListPayload struct {
	From         time.Time `json:"from"`
	To           time.Time `json:"to"`
	DashboardUID string    `json:"dashboard_uid,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	Limit        int       `json:"limit,omitempty"`
}

func (GrafanaAnnotationListPayload) Capability() string { return "annotation.list" }

func (p GrafanaAnnotationListPayload) Validate() error {
	if err := validateGrafanaWindow("annotation.list", p.From, p.To, GrafanaMaxAnnotationWindow); err != nil {
		return err
	}
	if p.DashboardUID != "" && !grafanaUID.MatchString(p.DashboardUID) {
		return errors.New("annotation.list dashboard_uid is invalid")
	}
	if err := validateGrafanaTags("annotation.list", p.Tags); err != nil {
		return err
	}
	return validateGrafanaLimit("annotation.list", p.Limit)
}

// GrafanaDatasourceQueryPayload runs datasource queries over a required,
// bounded time window. Queries are the datasource's own query objects (for
// example {"refId": "A", "expr": "up"}); MaxRows 0 means GrafanaMaxQueryRows.
type GrafanaDatasourceQueryPayload struct {
	From    time.Time        `json:"from"`
	To      time.Time        `json:"to"`
	MaxRows int              `json:"max_rows,omitempty"`
	Queries []map[string]any `json:"queries"`
}

func (GrafanaDatasourceQueryPayload) Capability() string { return "datasource.query" }

func (p GrafanaDatasourceQueryPayload) Validate() error {
	if err := validateGrafanaWindow("datasource.query", p.From, p.To, GrafanaMaxQueryWindow); err != nil {
		return err
	}
	if p.MaxRows < 0 || p.MaxRows > GrafanaMaxQueryRows {
		return fmt.Errorf("datasource.query max_rows must be between 1 and %d", GrafanaMaxQueryRows)
	}
	if len(p.Queries) == 0 || len(p.Queries) > GrafanaMaxQueries {
		return fmt.Errorf("datasource.query needs between 1 and %d queries", GrafanaMaxQueries)
	}
	for _, q := range p.Queries {
		if len(q) == 0 {
			return errors.New("datasource.query queries must not be empty")
		}
	}
	return nil
}

// EffectiveMaxRows is the row cap the runtime applies.
func (p GrafanaDatasourceQueryPayload) EffectiveMaxRows() int {
	if p.MaxRows == 0 {
		return GrafanaMaxQueryRows
	}
	return p.MaxRows
}

func validateGrafanaWindow(capability string, from, to time.Time, max time.Duration) error {
	if from.IsZero() || to.IsZero() || !from.Before(to) {
		return fmt.Errorf("%s needs from before to", capability)
	}
	if to.Sub(from) > max {
		return fmt.Errorf("%s window must be at most %s", capability, max)
	}
	return nil
}

func validateGrafanaTags(capability string, tags []string) error {
	if len(tags) > 10 {
		return fmt.Errorf("%s accepts at most 10 tags", capability)
	}
	for _, tag := range tags {
		if tag == "" || len(tag) > 128 {
			return fmt.Errorf("%s tags must be 1 to 128 characters", capability)
		}
	}
	return nil
}

func validateGrafanaLimit(capability string, limit int) error {
	if limit < 0 || limit > GrafanaMaxSearchLimit {
		return fmt.Errorf("%s limit must be between 1 and %d", capability, GrafanaMaxSearchLimit)
	}
	return nil
}
