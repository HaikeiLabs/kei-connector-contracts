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

package manifest

import (
	"errors"
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/validate"
)

// FileName is the generated file, relative to the plugin's out directory.
const FileName = "connector-manifest.v1.json"

// Options are the plugin parameters: contracts_version (required, the
// release tag the manifest ships in) and legacy (the path of
// legacy_capabilities.yaml, relative to the plugin's working directory).
type Options struct {
	Version    string
	LegacyPath string
}

// Set is a protogen.Options.ParamFunc.
func (o *Options) Set(name, value string) error {
	switch name {
	case "contracts_version":
		o.Version = value
	case "legacy":
		o.LegacyPath = value
	default:
		return fmt.Errorf("protoc-gen-kei-manifest: unknown parameter %q", name)
	}
	return nil
}

// Generate validates every operation and writes the manifest.
func Generate(gen *protogen.Plugin, opts Options) error {
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	if opts.Version == "" {
		return errors.New("protoc-gen-kei-manifest: contracts_version is required")
	}
	set, err := model.Load(gen.Files)
	if err != nil {
		return err
	}
	if errs := validate.Check(set); len(errs) > 0 {
		return errors.Join(errs...)
	}
	var legacy []LegacyRow
	if opts.LegacyPath != "" {
		b, err := os.ReadFile(opts.LegacyPath)
		if err != nil {
			return fmt.Errorf("protoc-gen-kei-manifest: %w", err)
		}
		if legacy, err = ParseLegacy(b); err != nil {
			return err
		}
	}
	m, err := Build(set, opts.Version, legacy)
	if err != nil {
		return err
	}
	b, err := Marshal(m)
	if err != nil {
		return err
	}
	_, err = gen.NewGeneratedFile(FileName, "").Write(b)
	return err
}
