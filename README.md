# SuperBotGo SDK

SDK repository for SuperBotGo plugins.

## Layout

- `protocol/` - versioned WASM plugin wire protocol, schemas, fixtures, and conformance tests.
- `go-sdk/` - Go SDK for building WASM plugins.

## Releases

Protocol and language SDKs are released with separate tag prefixes:

- `protocol/v4.0.0` for the wire contract.
- `go-sdk/v0.4.1` for the Go SDK module.

### Go SDK

Go SDK releases are created from GitHub Actions:

1. Open `Actions` -> `Release Go SDK`.
2. Run the workflow with `version`, for example `0.4.1`.
3. The workflow validates the protocol and Go SDK, creates tag `go-sdk/v0.4.1`, and publishes the GitHub Release.

Use the optional `prerelease` flag for versions like `0.4.1-rc.1`.

### Protocol

Protocol releases are still tag-based:

```bash
git tag protocol/v4.0.0
git push origin protocol/v4.0.0
```
