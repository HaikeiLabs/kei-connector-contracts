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

package render_test

import (
	"flag"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/plugintest"
	"github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei/internal/render"
)

var update = flag.Bool("update", false, "rewrite the golden files")

var goldenProtos = []string{
	"kei/testconn/v1/connector.proto",
	"kei/testconn/v1/things.proto",
	"kei/testsql/v1/connector.proto",
}

// generate runs the plugin over the golden protos and returns its files.
func generate(t *testing.T, parameter string) map[string]string {
	t.Helper()
	var opts render.Options
	gen := plugintest.Plugin(t, "testdata/protos", goldenProtos, parameter, opts.Set)
	if err := render.Generate(gen, opts); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	resp := gen.Response()
	if resp.Error != nil {
		t.Fatalf("plugin error: %s", resp.GetError())
	}
	out := map[string]string{}
	for _, f := range resp.File {
		out[f.GetName()] = f.GetContent()
	}
	return out
}

func checkGolden(t *testing.T, dir string, got map[string]string) {
	t.Helper()
	if *update {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
		for name, content := range got {
			p := filepath.Join(dir, name)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return
	}
	var want []string
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, p)
		want = append(want, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("reading goldens (run with -update): %v", err)
	}
	var names []string
	for n := range got {
		names = append(names, n)
	}
	sort.Strings(names)
	sort.Strings(want)
	if len(names) != len(want) {
		t.Fatalf("generated files %v, golden files %v", names, want)
	}
	for i, n := range names {
		if n != want[i] {
			t.Fatalf("generated files %v, golden files %v", names, want)
		}
		b, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != got[n] {
			t.Errorf("%s differs from its golden file; run go test ./internal/render -update and review the diff", n)
		}
	}
}

func TestGoGolden(t *testing.T) {
	checkGolden(t, "testdata/golden/go", generate(t, "paths=source_relative,lang=go"))
}

func TestTSGolden(t *testing.T) {
	checkGolden(t, "testdata/golden/ts", generate(t, "paths=source_relative,lang=ts"))
}

func TestGoGolden_Positional(t *testing.T) {
	checkGolden(t, "testdata/golden/go-positional", generate(t, "paths=source_relative,lang=go,signature=positional"))
}

func TestGenerate_Deterministic(t *testing.T) {
	a := generate(t, "paths=source_relative,lang=go")
	b := generate(t, "paths=source_relative,lang=go")
	for n := range a {
		if a[n] != b[n] {
			t.Fatalf("%s differs between two runs", n)
		}
	}
}

func TestOptions_Rejects(t *testing.T) {
	var o render.Options
	if err := o.Set("lang", "cobol"); err == nil {
		t.Errorf("lang=cobol: want an error")
	}
	if err := o.Set("signature", "named"); err == nil {
		t.Errorf("signature=named: want an error")
	}
	if err := o.Set("color", "blue"); err == nil {
		t.Errorf("unknown parameter: want an error")
	}
}
