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

// Package render emits the generated operations for one language from the
// shared model, with one text/template per language.
package render

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/model"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/validate"
)

// Options are the plugin parameters: lang=go|python|ts|rust and
// signature=message|positional.
type Options struct {
	Lang      string
	Signature string
}

// Set is a protogen.Options.ParamFunc.
func (o *Options) Set(name, value string) error {
	switch name {
	case "lang":
		switch value {
		case "go", "python", "ts", "rust":
			o.Lang = value
		default:
			return fmt.Errorf("protoc-gen-kei-ops: unknown lang %q (want go, python, ts or rust)", value)
		}
	case "signature":
		switch value {
		case "message", "positional":
			o.Signature = value
		default:
			return fmt.Errorf("protoc-gen-kei-ops: unknown signature %q (want message or positional)", value)
		}
	default:
		return fmt.Errorf("protoc-gen-kei-ops: unknown parameter %q", name)
	}
	return nil
}

func (o Options) positional() bool { return o.Signature == "positional" }

// Generate validates every operation in gen's files and emits the
// operations of each connector declared in a file being generated.
func Generate(gen *protogen.Plugin, opts Options) error {
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	if opts.Lang == "" {
		opts.Lang = "go"
	}
	set, err := model.Load(gen.Files)
	if err != nil {
		return err
	}
	if errs := validate.Check(set); len(errs) > 0 {
		return errors.Join(errs...)
	}
	for _, c := range set.Connectors {
		if !c.File.Generate || len(c.Ops) == 0 {
			continue
		}
		switch opts.Lang {
		case "go":
			err = renderGo(gen, c, opts)
		case "ts":
			err = renderTS(gen, c, opts)
		case "rust":
			err = renderRust(gen, c, opts)
		case "python":
			err = renderPython(gen, c, opts)
		default:
			err = fmt.Errorf("protoc-gen-kei-ops: lang %q is not implemented yet", opts.Lang)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
