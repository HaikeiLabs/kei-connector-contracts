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

// Package plugintest compiles .proto files in process and builds the
// protogen.Plugin a plugin would receive from buf, for tests.
package plugintest

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// RepoProto is the repository's proto/ directory, so test protos can import
// kei/contracts/v1/options.proto.
func RepoProto() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "proto")
}

// Plugin compiles files (relative to importDir) and returns the plugin that
// generates them, with the given parameter string (for example
// "paths=source_relative,lang=go"). Unknown parameters are passed to
// paramFunc, which may be nil.
func Plugin(t *testing.T, importDir string, files []string, parameter string, paramFunc func(name, value string) error) *protogen.Plugin {
	t.Helper()
	if paramFunc == nil {
		paramFunc = func(string, string) error { return nil }
	}
	gen, err := protogen.Options{ParamFunc: paramFunc}.New(Request(t, importDir, files, parameter))
	if err != nil {
		t.Fatalf("protogen: %v", err)
	}
	return gen
}

// Request compiles files and builds the CodeGeneratorRequest. The request is
// round-tripped through the wire format, so the kei.contracts.v1 extensions
// resolve to the linked Go types exactly as in a real plugin run.
func Request(t *testing.T, importDir string, files []string, parameter string) *pluginpb.CodeGeneratorRequest {
	t.Helper()
	c := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: []string{importDir, RepoProto()},
		}),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	compiled, err := c.Compile(context.Background(), files...)
	if err != nil {
		t.Fatalf("compile %v: %v", files, err)
	}
	var protos []*descriptorpb.FileDescriptorProto
	seen := map[string]bool{}
	for _, f := range compiled {
		collect(f, seen, &protos)
	}
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: files,
		Parameter:      proto.String(parameter),
		ProtoFile:      protos,
	}
	b, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	out := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(b, out); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	return out
}

// collect appends fd's dependencies, then fd, in topological order.
func collect(fd protoreflect.FileDescriptor, seen map[string]bool, out *[]*descriptorpb.FileDescriptorProto) {
	if seen[fd.Path()] {
		return
	}
	seen[fd.Path()] = true
	imports := fd.Imports()
	for i := 0; i < imports.Len(); i++ {
		collect(imports.Get(i).FileDescriptor, seen, out)
	}
	*out = append(*out, protodesc.ToFileDescriptorProto(fd))
}
