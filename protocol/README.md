# SuperBotGo Plugin Protocol

Protocol versions live in separate subdirectories:

- `v4/` - current WASM plugin protocol.

Each version directory owns its `protocol.md`, JSON Schemas, fixtures, and
conformance tests. New breaking wire-contract changes should create a new
version directory instead of modifying older protocol versions in place.

In a combined SDK repository, release protocol and language SDKs with separate
tag prefixes, for example `protocol/v4.0.0` for the wire contract and
`go-sdk/v0.4.1` for the Go SDK module.
