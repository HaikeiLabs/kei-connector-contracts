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
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts"
)

type S3Object struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	ETag         string    `json:"etag"`
	LastModified time.Time `json:"last_modified"`
}

// ObjectListPayload requests object.list. The list scope is carried by the
// invocation resource (s3://<bucket> or s3://<bucket>/<prefix>), and Prefix
// must agree with it.
type ObjectListPayload struct {
	Prefix string `json:"prefix"`
}

func (ObjectListPayload) Capability() string { return "object.list" }

// ObjectReadPayload requests object.read. The key is carried by the
// invocation resource (s3://<bucket>/<key>), and Key must agree with it.
type ObjectReadPayload struct {
	Key string `json:"key"`
}

func (ObjectReadPayload) Capability() string { return "object.read" }

// S3Backend is the seam for S3 reads.
type S3Backend interface {
	ListObjects(ctx context.Context, bucket, prefix string) ([]S3Object, error)
	GetObject(ctx context.Context, bucket, key string) ([]byte, error)
}

// MemoryS3 is the in-memory S3 backend used until a real one exists. Objects
// are keyed by bucket/key.
type MemoryS3 struct {
	Objects  map[string]S3Object
	Contents map[string][]byte
}

func s3Key(bucket, key string) string { return bucket + "/" + key }

func (m MemoryS3) ListObjects(_ context.Context, bucket, prefix string) ([]S3Object, error) {
	var out []S3Object
	for k, o := range m.Objects {
		if !strings.HasPrefix(k, bucket+"/") {
			continue
		}
		key := strings.TrimPrefix(k, bucket+"/")
		if prefix == "" || strings.HasPrefix(key, prefix) {
			out = append(out, o)
		}
	}
	return out, nil
}

func (m MemoryS3) GetObject(_ context.Context, bucket, key string) ([]byte, error) {
	data, ok := m.Contents[s3Key(bucket, key)]
	if !ok {
		return nil, errors.New("object not found")
	}
	return data, nil
}

// S3Client serves object.list and object.read, scoped to the connector's
// allowed resources (buckets) and allowed prefixes (key prefixes).
type S3Client struct {
	backend S3Backend
}

func NewS3(backend S3Backend) *S3Client { return &S3Client{backend: backend} }

func (c *S3Client) Provider() connectors.Provider { return connectors.ProviderS3 }

func (c *S3Client) Invoke(ctx context.Context, meta connectors.Metadata, inv connectors.Invocation, payload Payload) (Result, error) {
	if err := checkProvider(meta, c.Provider()); err != nil {
		return Result{}, err
	}
	if err := Guard(meta, inv); err != nil {
		return Result{}, err
	}
	if err := matchCapability(inv, payload); err != nil {
		return Result{}, err
	}
	bucket, key, ok := s3Resource(inv.Resource)
	if !ok {
		return Result{}, errors.New("resource must be s3://<bucket> or s3://<bucket>/<key>")
	}
	switch p := payload.(type) {
	case ObjectListPayload:
		if p.Prefix != key {
			return Result{}, errors.New("payload prefix does not match the invocation resource")
		}
		objects, err := c.backend.ListObjects(ctx, bucket, p.Prefix)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: objects}, nil
	case ObjectReadPayload:
		if key == "" {
			return Result{}, errors.New("object.read requires a key in the resource")
		}
		if p.Key != key {
			return Result{}, errors.New("payload key does not match the invocation resource")
		}
		data, err := c.backend.GetObject(ctx, bucket, key)
		if err != nil {
			return Result{}, err
		}
		return Result{Capability: inv.Capability, Data: data}, nil
	default:
		return Result{}, fmt.Errorf("unsupported payload %T", payload)
	}
}

func s3Resource(resource string) (bucket, key string, ok bool) {
	const scheme = "s3://"
	if !strings.HasPrefix(resource, scheme) {
		return "", "", false
	}
	rest := strings.TrimPrefix(resource, scheme)
	if rest == "" {
		return "", "", false
	}
	if i := strings.IndexByte(rest, '/'); i >= 0 {
		return rest[:i], rest[i+1:], true
	}
	return rest, "", true
}
