module github.com/HaikeiLabs/kei-connector-contracts/tools/protoc-gen-kei

go 1.27

require (
	github.com/HaikeiLabs/kei-connector-contracts/gen/go v0.0.0
	github.com/bufbuild/protocompile v0.14.1
	google.golang.org/protobuf v1.36.10
)

require golang.org/x/sync v0.8.0 // indirect

// The plugins read the kei.contracts.v1 options through the generated Go
// types in this repository, never through a published version.
replace github.com/HaikeiLabs/kei-connector-contracts/gen/go => ../../gen/go
