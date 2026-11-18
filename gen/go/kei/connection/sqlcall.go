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
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

// DefaultMaxRows applies when a query request leaves max_rows unset.
const DefaultMaxRows = 1000

// maxExactInt is 2^53: larger integers are returned as strings, because a
// google.protobuf.Value number is a double.
const maxExactInt = 1 << 53

// CallSQL runs a query operation. The request must have the fields sql
// (string), params (repeated google.protobuf.Value) and max_rows (optional
// int64); the response columns (repeated {name, type}), rows (repeated
// google.protobuf.ListValue) and truncated (bool). The generator validates
// that shape.
//
// Cell values follow the database contract: NULL is null, an integer beyond
// ±2^53 and a []byte value (NUMERIC, DECIMAL, text) are strings, a time is
// an RFC 3339 string in UTC.
func CallSQL(ctx context.Context, conn SQL, op Op, req, resp proto.Message) error {
	if conn == nil {
		return fmt.Errorf("connection: %s: nil connection", op.ID)
	}
	rm := req.ProtoReflect()
	rf := rm.Descriptor().Fields()
	sqlFd, paramsFd, maxFd := rf.ByName("sql"), rf.ByName("params"), rf.ByName("max_rows")
	if sqlFd == nil || paramsFd == nil || maxFd == nil {
		return fmt.Errorf("connection: %s: request is not a query request", op.ID)
	}
	var args []any
	params := rm.Get(paramsFd).List()
	for i := 0; i < params.Len(); i++ {
		v, ok := params.Get(i).Message().Interface().(*structpb.Value)
		if !ok {
			return fmt.Errorf("connection: %s: params must be google.protobuf.Value", op.ID)
		}
		args = append(args, v.AsInterface())
	}
	maxRows := int64(DefaultMaxRows)
	if rm.Has(maxFd) && rm.Get(maxFd).Int() > 0 {
		maxRows = rm.Get(maxFd).Int()
	}
	cols, rows, truncated, err := conn.Query(ctx, rm.Get(sqlFd).String(), args, maxRows)
	if err != nil {
		return fmt.Errorf("connection: %s: %w", op.ID, err)
	}

	pm := resp.ProtoReflect()
	pf := pm.Descriptor().Fields()
	colsFd, rowsFd, truncFd := pf.ByName("columns"), pf.ByName("rows"), pf.ByName("truncated")
	if colsFd == nil || rowsFd == nil || truncFd == nil {
		return fmt.Errorf("connection: %s: response is not a query response", op.ID)
	}
	colList := pm.Mutable(colsFd).List()
	for _, c := range cols {
		e := colList.NewElement()
		ef := e.Message().Descriptor().Fields()
		e.Message().Set(ef.ByName("name"), protoreflect.ValueOfString(c.Name))
		e.Message().Set(ef.ByName("type"), protoreflect.ValueOfString(strings.ToLower(c.Type)))
		colList.Append(e)
	}
	rowList := pm.Mutable(rowsFd).List()
	for _, row := range rows {
		lv := &structpb.ListValue{}
		for _, cell := range row {
			lv.Values = append(lv.Values, ToValue(cell))
		}
		rowList.Append(protoreflect.ValueOfMessage(lv.ProtoReflect()))
	}
	pm.Set(truncFd, protoreflect.ValueOfBool(truncated))
	return nil
}

// ToValue converts one database cell to a google.protobuf.Value.
func ToValue(v any) *structpb.Value {
	switch x := v.(type) {
	case nil:
		return structpb.NewNullValue()
	case bool:
		return structpb.NewBoolValue(x)
	case int64:
		return intValue(x)
	case int:
		return intValue(int64(x))
	case int32:
		return intValue(int64(x))
	case int16:
		return intValue(int64(x))
	case int8:
		return intValue(int64(x))
	case uint64:
		if x > maxExactInt {
			return structpb.NewStringValue(strconv.FormatUint(x, 10))
		}
		return structpb.NewNumberValue(float64(x))
	case uint32:
		return structpb.NewNumberValue(float64(x))
	case float64:
		return floatValue(x)
	case float32:
		return floatValue(float64(x))
	case string:
		return structpb.NewStringValue(x)
	case []byte:
		return structpb.NewStringValue(string(x))
	case time.Time:
		return structpb.NewStringValue(x.UTC().Format(time.RFC3339Nano))
	default:
		return structpb.NewStringValue(fmt.Sprint(x))
	}
}

func intValue(n int64) *structpb.Value {
	if n > maxExactInt || n < -maxExactInt {
		return structpb.NewStringValue(strconv.FormatInt(n, 10))
	}
	return structpb.NewNumberValue(float64(n))
}

func floatValue(f float64) *structpb.Value {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return structpb.NewStringValue(strconv.FormatFloat(f, 'g', -1, 64))
	}
	return structpb.NewNumberValue(f)
}
