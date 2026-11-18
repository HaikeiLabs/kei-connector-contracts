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

package connection

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Canonical request JSON, shared by every language's runtime so the
// cross-language parity test can compare request bodies byte for byte:
//
//   - object keys sorted by byte order, no insignificant whitespace;
//   - integers (including 64-bit) as exact JSON numbers, never strings;
//   - strings with only the escapes JSON requires (no HTML escaping, UTF-8
//     kept as is);
//   - enums by name, bytes as standard base64;
//   - only populated fields; well-known types in their protobuf JSON form.

// WireName is the key a field uses on the wire: its json_name when that
// differs from the protobuf default (lowerCamelCase), otherwise the proto
// field name. Providers mostly use snake_case keys, which are proto names;
// a field whose provider key differs declares json_name.
func WireName(fd protoreflect.FieldDescriptor) string {
	if fd.JSONName() != defaultJSONName(string(fd.Name())) {
		return fd.JSONName()
	}
	return string(fd.Name())
}

// nameFunc names a field inside an encoded message: WireName in a REST body,
// the JSON name (camelCase) in GraphQL variables.
type nameFunc func(protoreflect.FieldDescriptor) string

func jsonName(fd protoreflect.FieldDescriptor) string { return fd.JSONName() }

// defaultJSONName is protoc's lowerCamelCase JSON name for a field name.
func defaultJSONName(name string) string {
	var b strings.Builder
	upper := false
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c == '_':
			upper = true
		case upper && 'a' <= c && c <= 'z':
			b.WriteByte(c - 'a' + 'A')
			upper = false
		default:
			b.WriteByte(c)
			upper = false
		}
	}
	return b.String()
}

// object is a JSON object under construction.
type object map[string]json.RawMessage

func (o object) encode() []byte {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		b.Write(jsonString(k))
		b.WriteByte(':')
		b.Write(o[k])
	}
	b.WriteByte('}')
	return b.Bytes()
}

// jsonString encodes s without HTML escaping.
func jsonString(s string) []byte {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s) // a string never fails to encode
	return bytes.TrimRight(b.Bytes(), "\n")
}

// encodeMessage encodes every populated field of m.
func encodeMessage(m protoreflect.Message, name nameFunc) (json.RawMessage, error) {
	if isWellKnown(m.Descriptor()) {
		return encodeWellKnown(m)
	}
	o := object{}
	var err error
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		var raw json.RawMessage
		raw, err = encodeField(fd, v, name)
		if err != nil {
			return false
		}
		o[name(fd)] = raw
		return true
	})
	if err != nil {
		return nil, err
	}
	return o.encode(), nil
}

func encodeField(fd protoreflect.FieldDescriptor, v protoreflect.Value, name nameFunc) (json.RawMessage, error) {
	switch {
	case fd.IsList():
		l := v.List()
		var b bytes.Buffer
		b.WriteByte('[')
		for i := 0; i < l.Len(); i++ {
			if i > 0 {
				b.WriteByte(',')
			}
			raw, err := encodeSingular(fd, l.Get(i), name)
			if err != nil {
				return nil, err
			}
			b.Write(raw)
		}
		b.WriteByte(']')
		return b.Bytes(), nil
	case fd.IsMap():
		o := object{}
		var err error
		v.Map().Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
			var raw json.RawMessage
			raw, err = encodeSingular(fd.MapValue(), mv, name)
			o[k.String()] = raw
			return err == nil
		})
		if err != nil {
			return nil, err
		}
		return o.encode(), nil
	default:
		return encodeSingular(fd, v, name)
	}
}

func encodeSingular(fd protoreflect.FieldDescriptor, v protoreflect.Value, name nameFunc) (json.RawMessage, error) {
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return encodeMessage(v.Message(), name)
	case protoreflect.EnumKind:
		if ev := fd.Enum().Values().ByNumber(v.Enum()); ev != nil {
			return jsonString(string(ev.Name())), nil
		}
		return json.RawMessage(strconv.FormatInt(int64(v.Enum()), 10)), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		f := v.Float()
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil, fmt.Errorf("field %s: %v is not valid JSON", fd.Name(), f)
		}
		b, _ := json.Marshal(f)
		return b, nil
	default:
		s, err := scalarString(fd, v)
		if err != nil {
			return nil, err
		}
		switch fd.Kind() {
		case protoreflect.StringKind, protoreflect.BytesKind:
			return jsonString(s), nil
		default: // bool and integers are bare literals
			return json.RawMessage(s), nil
		}
	}
}

// scalarString formats a scalar for a path, a query value or a JSON literal.
func scalarString(fd protoreflect.FieldDescriptor, v protoreflect.Value) (string, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return strconv.FormatBool(v.Bool()), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return strconv.FormatInt(v.Int(), 10), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return strconv.FormatUint(v.Uint(), 10), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return strconv.FormatFloat(v.Float(), 'g', -1, 64), nil
	case protoreflect.StringKind:
		return v.String(), nil
	case protoreflect.BytesKind:
		return base64.StdEncoding.EncodeToString(v.Bytes()), nil
	case protoreflect.EnumKind:
		if ev := fd.Enum().Values().ByNumber(v.Enum()); ev != nil {
			return string(ev.Name()), nil
		}
		return strconv.FormatInt(int64(v.Enum()), 10), nil
	default:
		return "", fmt.Errorf("field %s: a %s cannot be a path or query value", fd.Name(), fd.Kind())
	}
}

func isWellKnown(md protoreflect.MessageDescriptor) bool {
	return md.ParentFile() != nil && md.ParentFile().Package() == "google.protobuf"
}

// encodeWellKnown re-encodes protojson's output canonically (protojson's
// whitespace is deliberately unstable).
func encodeWellKnown(m protoreflect.Message) (json.RawMessage, error) {
	b, err := protojson.Marshal(m.Interface())
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return canonical(v), nil
}

// canonical encodes a decoded JSON value (with json.Number) canonically.
func canonical(v any) json.RawMessage {
	switch x := v.(type) {
	case nil:
		return json.RawMessage("null")
	case bool:
		return json.RawMessage(strconv.FormatBool(x))
	case json.Number:
		return json.RawMessage(x.String())
	case string:
		return jsonString(x)
	case []any:
		var b bytes.Buffer
		b.WriteByte('[')
		for i, e := range x {
			if i > 0 {
				b.WriteByte(',')
			}
			b.Write(canonical(e))
		}
		b.WriteByte(']')
		return b.Bytes()
	case map[string]any:
		o := object{}
		for k, e := range x {
			o[k] = canonical(e)
		}
		return o.encode()
	default:
		return jsonString(fmt.Sprint(x))
	}
}
