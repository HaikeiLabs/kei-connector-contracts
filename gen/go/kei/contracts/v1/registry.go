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

package contracts

import (
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// All returns every operation declared by a request message in a file
// descriptor linked into this binary, sorted by id.
//
// A proto file enters protoregistry.GlobalFiles when its generated Go
// package is imported: a program that imports only this package sees no
// operations, a program that also imports a connector package sees that
// connector's operations. Discovery is therefore static — exactly what a
// policy or tooling layer wants.
func All() []*Operation {
	var ops []*Operation
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		collectOperations(fd.Messages(), &ops)
		return true
	})
	slices.SortFunc(ops, func(a, b *Operation) int {
		return strings.Compare(a.GetId(), b.GetId())
	})
	return ops
}

// Lookup returns the operation with the given id, if one is registered.
func Lookup(id string) (*Operation, bool) {
	for _, op := range All() {
		if op.GetId() == id {
			return op, true
		}
	}
	return nil, false
}

// collectOperations appends the operation declared by each message (nested
// messages included) that carries the operation option.
func collectOperations(messages protoreflect.MessageDescriptors, ops *[]*Operation) {
	for i := 0; i < messages.Len(); i++ {
		md := messages.Get(i)
		if op := operationOf(md); op != nil {
			*ops = append(*ops, op)
		}
		collectOperations(md.Messages(), ops)
	}
}

// operationOf reads the kei.contracts.v1.operation extension from a
// message's options; unset means no operation.
func operationOf(md protoreflect.MessageDescriptor) *Operation {
	opts, ok := md.Options().(*descriptorpb.MessageOptions)
	if !ok || opts == nil {
		return nil
	}
	if !proto.HasExtension(opts, E_Operation) {
		return nil
	}
	op, ok := proto.GetExtension(opts, E_Operation).(*Operation)
	if !ok {
		return nil
	}
	return op
}
