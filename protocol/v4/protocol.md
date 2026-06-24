# SuperBotGo Plugin Protocol v4

This document describes the wire contract between the SuperBotGo host and a
WASM plugin. Language SDKs are convenience layers over this protocol; they are
not the protocol itself.

## Contract Boundaries

The shared contract for WASM plugins lives in this directory as documentation,
schemas, fixtures, and conformance tests. The host mirrors these wire DTOs in
`internal/wasm/protocol`.

The Go-native plugin interface is a separate host contract in
`internal/plugin/contract`. It may use internal domain types such as user IDs,
channel types, option maps, and file references. The WASM adapter converts
between `internal/plugin/contract` and `internal/wasm/protocol` explicitly, so a
future SDK in another language depends only on this protocol directory and not
on Go-native host types.

## Versioning

Current protocol version: `4`.

The current Go SDK exposes this as `ProtocolVersion` and serializes it in plugin
metadata as `sdk_version`. That field name is historical and means protocol
version for compatibility with existing plugins.

Compatibility rules:

- A host may run plugins whose protocol version is less than or equal to the
  host's maximum supported protocol version.
- A plugin must not require a protocol version greater than the host supports.
- Minor protocol extensions should add optional fields.
- Removing fields, changing field meaning, or changing encodings requires a new
  protocol version.

## Execution Model

The host executes the plugin as a WASI command module for one action at a time.

For each action:

1. The host instantiates the compiled WASM module.
2. The host sets `PLUGIN_ACTION`.
3. The host writes the action request to stdin as JSON, when the action has a
   request body.
4. The plugin writes the action response to stdout as JSON, when the action has
   a response body.
5. The module exits with code `0` for a handled action.

The host may set `PLUGIN_CONFIG` to the persisted plugin configuration encoded
as JSON for actions that need runtime config.

## Actions

### `meta`

Request body: empty.

Response body: `schemas/plugin-meta.schema.json`.

The response declares plugin identity, triggers, requirements, RPC methods,
configuration schema, and migrations.

### `configure`

Request body: plugin configuration JSON.

Response body:

- empty on success;
- `schemas/event-response.schema.json` with `error` on failure.

### `reconfigure`

Request body: `schemas/reconfigure-request.schema.json`.

Response body:

- empty on success;
- `schemas/event-response.schema.json` with `error` on failure.

### `handle_event`

Request body: `schemas/event-request.schema.json`.

Response body: `schemas/event-response.schema.json`.

Messenger replies are returned in `reply_blocks`. HTTP trigger responses are
returned in `data` as an object with `status_code`, optional `headers`, and
`body`.

### `handle_rpc`

Request body: `schemas/rpc-request.schema.json`.

Response body: `schemas/rpc-response.schema.json`.

`params` and `result` are byte arrays serialized by JSON as base64 strings.
The bytes currently contain msgpack payloads.

### `step_callback`

Request body: `schemas/step-callback-request.schema.json`.

Response body: `schemas/step-callback-response.schema.json`.

Callbacks are used by interactive messenger nodes for dynamic options,
pagination, validation, and conditional visibility.

### `migrate`

Request body: `schemas/migrate-request.schema.json`.

Response body: `schemas/migrate-response.schema.json`.

This action handles plugin data migration between plugin versions.

## Host Calls

Plugins can import host functions from the `env` module. Every host call uses
the same low-level ABI:

```text
env.<function>(ptr: i32, len: i32) -> i64
```

The request payload is msgpack bytes at `(ptr, len)` in plugin memory. When the
host returns bytes, the result packs the response pointer and length:

```text
result = (response_ptr << 32) | response_len
```

The plugin must export:

```text
alloc(size: i32) -> i32
alloc_reset() -> nil
```

The host uses `alloc` to copy response bytes into plugin memory. A result value
of `0` means no response body.

Current host imports include:

- `http_request`
- `call_plugin`
- `publish_event`
- `kv_get`, `kv_set`, `kv_delete`, `kv_list`
- `notify_user`, `notify_chat`, `notify_users`, `notify_teacher`,
  `notify_students`
- `sql_open`, `sql_close`, `sql_exec`, `sql_query`, `sql_next`,
  `sql_rows_close`, `sql_begin`, `sql_end`
- `file_meta`, `file_read`, `file_read_into`, `file_url`, `file_store`
- `user_info`, `users_info`

### `user_info`

Requires plugin requirement type `user_info`.

Request (msgpack):

```
{ "user_id": int64 }
```

Response (msgpack):

```
{
  "id": int64,
  "full_name": string,
  "external_id": string,
  "is_teacher": bool
}
```

Returns basic information about a global user by their internal ID. `full_name`
is derived from the university `persons` record (last name + first name +
middle name) when available, otherwise falls back to the messenger username.
`external_id` is the university student/employee ID. `is_teacher` is `true`
when the user has an active record in `teacher_positions`. Returns an error if
the user is not found.

### `users_info`

Requires plugin requirement type `user_info`.

Request (msgpack):

```
{ "user_ids": []int64 }
```

Response (msgpack):

```
{
  "users": [
    {
      "id": int64,
      "full_name": string,
      "external_id": string,
      "is_teacher": bool,
      "positions": [
        {
          "position_type": string,
          "status": string,
          "nationality_type": string,
          "funding_type": string,
          "education_form": string,
          "group_code": string,
          "group_name": string,
          "program_name": string,
          "stream_name": string
        }
      ]
    }
  ]
}
```

Returns information for multiple users in a single call. Each entry includes the
same fields as `user_info` plus a `positions` list. Currently only student
positions (`position_type = "student"`) are returned, one entry per active
`student_positions` record. A user with no positions receives an empty
`positions` array. Users not found in the database are omitted from the result.

Host-call payload schemas are not yet the source of truth in this directory;
this first contract layer covers the JSON lifecycle messages. The host-call ABI
above is part of the protocol and should be expanded with msgpack payload
schemas before a second language SDK becomes public.

## Schemas

Schemas are written as JSON Schema Draft 2020-12 and live in `schemas/`.

The conformance test compiles every schema and validates fixtures under
`testdata/valid` and `testdata/invalid`.
