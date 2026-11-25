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

package validate_test

import (
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/plugintest"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/validate"
)

// TestValidate_Rules compiles one invalid proto per rule (section 5.1 of
// proposal 0002, plus the semantics and GraphQL rules) and expects exactly
// one error, naming that rule.
func TestValidate_Rules(t *testing.T) {
	rules := []string{
		validate.RuleDuplicateID,
		validate.RuleConnectorMismatch,
		validate.RuleWriteWithoutScope,
		validate.RuleUnknownPathVar,
		validate.RuleUndeclaredConnection,
		validate.RuleMissingResponse,
		validate.RuleMissingSemantics,
		validate.RuleSemanticsClassMismatch,
		validate.RuleGraphQLWithoutDocument,
		validate.RuleUnsupportedParam,
		validate.RuleMissingResultFields,
		validate.RuleUnsupportedResultField,
		validate.RuleResultFieldsOnRead,
	}
	for _, rule := range rules {
		t.Run(rule, func(t *testing.T) {
			gen := plugintest.Plugin(t, "testdata/invalid", []string{rule + ".proto"}, "", nil)
			set, err := model.Load(gen.Files)
			if err != nil {
				t.Fatalf("model.Load: %v", err)
			}
			errs := validate.Check(set)
			if len(errs) != 1 {
				t.Fatalf("Check = %d errors %v, want exactly 1 naming %s", len(errs), errs, rule)
			}
			if !strings.HasPrefix(errs[0].Error(), rule+":") {
				t.Fatalf("Check = %q, want the %s rule", errs[0], rule)
			}
		})
	}
}

// TestValidate_ValidSet checks that the golden test protos pass.
func TestValidate_ValidSet(t *testing.T) {
	gen := plugintest.Plugin(t, "../render/testdata/protos", []string{
		"kei/testconn/v1/connector.proto",
		"kei/testconn/v1/things.proto",
		"kei/testsql/v1/connector.proto",
	}, "", nil)
	set, err := model.Load(gen.Files)
	if err != nil {
		t.Fatalf("model.Load: %v", err)
	}
	if errs := validate.Check(set); len(errs) != 0 {
		t.Fatalf("Check = %v, want none", errs)
	}
	if len(set.Connectors) != 2 || len(set.Ops) == 0 {
		t.Fatalf("Load = %d connectors, %d ops", len(set.Connectors), len(set.Ops))
	}
}
