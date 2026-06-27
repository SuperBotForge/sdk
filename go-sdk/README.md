# SuperBotForge Go SDK

Go SDK for writing WASM plugins for the SuperBotGo platform.

## Quick start

```go
package main

import wasmplugin "github.com/SuperBotForge/sdk/go-sdk"

func main() {
    wasmplugin.Run(wasmplugin.Plugin{
        ID:      "my-plugin",
        Name:    "My Plugin",
        Version: "1.0.0",
        Requirements: []wasmplugin.Requirement{
            wasmplugin.NotifyReq("Send notifications to users").Build(),
            wasmplugin.KV("Store plugin state").Build(),
        },
        Triggers: []wasmplugin.Trigger{
            {
                Name: "/hello",
                Type: wasmplugin.TriggerMessenger,
                Descriptions: map[string]string{
                    "ru": "Поздороваться",
                    "en": "Say hello",
                },
                Handler: handleHello,
            },
        },
    })
}

func handleHello(ctx *wasmplugin.EventContext) error {
    ctx.Reply(wasmplugin.NewMessage("Hello, " + ctx.Param("name") + "!"))
    return nil
}
```

Build for WASM:

```sh
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm .
```

---

## EventContext

`EventContext` is passed to every handler. The active sub-struct tells you which trigger fired:

| Field | Not nil when |
|---|---|
| `ctx.Messenger` | Messenger command (`TriggerMessenger`) |
| `ctx.HTTP` | HTTP endpoint (`TriggerHTTP`) |
| `ctx.Cron` | Cron schedule (`TriggerCron`) |
| `ctx.Event` | Event bus topic (`TriggerEvent`) |

### Messenger fields

```go
ctx.Messenger.UserID      int64   // Global user ID
ctx.Messenger.ChannelType string  // "telegram", "vk", "discord", "mattermost"
ctx.Messenger.ChatID      string  // Platform chat ID
ctx.Messenger.ChatGroupID string  // Cross-messenger group ID (if any)
ctx.Messenger.CommandName string  // e.g. "/start"
ctx.Messenger.Params      map[string]string
ctx.Messenger.Locale      string  // "ru", "en", …
ctx.Messenger.Files       []FileRef
```

### HTTP fields

```go
ctx.HTTP.Method     string
ctx.HTTP.Path       string
ctx.HTTP.Query      map[string]string
ctx.HTTP.Headers    map[string]string
ctx.HTTP.Body       string
ctx.HTTP.RemoteAddr string
ctx.HTTP.Auth       *HTTPAuthInfo // non-nil when request is authenticated
```

### Cron fields

```go
ctx.Cron.ScheduleName string
ctx.Cron.FireTime     int64 // Unix timestamp
```

### Event bus fields

```go
ctx.Event.Topic   string
ctx.Event.Payload []byte
ctx.Event.Source  string
```

### Helpers

```go
ctx.Param("key")   // shortcut for ctx.Messenger.Params["key"]
ctx.Locale()       // shortcut for ctx.Messenger.Locale
ctx.Config("key", "fallback")  // plugin config value
ctx.Log("message")
ctx.LogError("message")
```

---

## Responding

### Messenger commands

```go
ctx.Reply(msg Message)
```

### HTTP handlers

```go
// Low-level
ctx.SetHTTPResponse(statusCode int, headers map[string]string, body string)

// Convenience: sets Content-Type: application/json
ctx.JSON(statusCode int, v interface{})
```

---

## Message builder

`Message` is a block-based rich message. Blocks are rendered natively on each platform (Telegram inline keyboard, Discord embeds, etc.).

```go
// Plain text
wasmplugin.NewMessage("Готово!")

// Localized text (host picks the right locale per recipient)
wasmplugin.NewLocalizedMessage(map[string]string{"ru": "Готово!", "en": "Done!"})
```

Chain methods to add blocks:

```go
msg := wasmplugin.NewMessage("Расписание обновлено").
    StyledText("Изменения вступят в силу завтра", wasmplugin.StyleQuote).
    Mention(platformUserID).
    Link("https://schedule.example.com", "Открыть расписание").
    File(ref, "schedule.pdf").
    Options("Что делать?",
        wasmplugin.Opt("📊 Мой статус", "/queue_status"),
        wasmplugin.Opt("🏠 Главное меню", "/start"),
    )
```

| Method | Block type | Notes |
|---|---|---|
| `.Text(text)` | Plain text | Appends a text block |
| `.StyledText(text, style)` | Styled text | Styles: `StylePlain`, `StyleHeader`, `StyleSubheader`, `StyleCode`, `StyleQuote` |
| `.LocalizedText(map)` | Localized text | Host resolves locale per recipient |
| `.LocalizedStyledText(map, style)` | Localized styled text | |
| `.Mention(platformUserID)` | User mention | `platformUserID` is the channel-specific user ID (e.g. Telegram user ID as string) |
| `.File(ref, caption)` | File attachment | `ref` is a `FileRef` from file store or incoming event |
| `.Link(url, label)` | Hyperlink | |
| `.Image(url)` | Inline image | |
| `.Options(prompt, opts...)` | Inline keyboard | Button click sends `Option.Value` as input; `/command` values trigger commands |

Create an option:

```go
wasmplugin.Opt("Button label", "/command_or_value")
```

---

## User API

Requires `wasmplugin.UserInfoReq("...").Build()` in `Requirements`.

### GetUserInfo

```go
info, err := wasmplugin.GetUserInfo(userID int64) (*UserInfo, error)
// or
info, err := ctx.GetUserInfo(userID)
```

```go
type UserInfo struct {
    ID            int64
    FullName      string
    ExternalID    string  // university person external_id (matches persons.external_id in DB)
    TsuAccountsID string  // TSU accounts.tsu.ru account ID
    TsuLinked     bool    // true if user has linked a TSU account
    IsTeacher     bool
    IsStudent     bool
    IsDeanOffice  bool
}
```

### GetUsersInfo

Fetches multiple users at once, including their university positions.

```go
users, err := wasmplugin.GetUsersInfo(userIDs []int64) ([]UserInfoFull, error)
// or
users, err := ctx.GetUsersInfo(userIDs)
```

```go
type UserInfoFull struct {
    UserInfo
    Positions []UserPosition
}

type UserPosition struct {
    PositionType    string // "student" or "teacher"
    Status          string
    NationalityType string
    FundingType     string
    EducationForm   string
    FacultyName     string
    DepartmentName  string
    ProgramName     string
    StreamName      string
    GroupCode       string
    GroupName       string
}
```

### ListUsers

Paginated list of all users with positions.

```go
users, total, err := wasmplugin.ListUsers(page, pageSize int) ([]UserInfoFull, int, error)
// or
users, total, err := ctx.ListUsers(page, pageSize)
```

Page is 0-indexed. Example: `ListUsers(0, 50)` returns the first 50 users.

---

## KV Store

Plugin-scoped key-value store. All keys are namespaced to the plugin automatically.
Requires `wasmplugin.KV("...").Build()` in `Requirements`.

```go
// Get a value
value, found, err := ctx.KVGet(key string) (string, bool, error)

// Set a value (persists until plugin is unloaded or host restarts)
err := ctx.KVSet(key, value string) error

// Set with TTL (auto-expires)
err := ctx.KVSetWithTTL(key, value string, ttl time.Duration) error

// Delete
err := ctx.KVDelete(key string) error

// List keys by prefix (pass "" to list all)
keys, err := ctx.KVList(prefix string) ([]string, error)
```

The same methods are available on `MigrateContext` for data migration in the `Migrate` callback.

---

## Notifications

Requires `wasmplugin.NotifyReq("...").Build()` in `Requirements`.

Priority levels:

| Constant | Value | Behaviour |
|---|---|---|
| `PriorityLow` | 0 | Delayed outside work hours, no sound |
| `PriorityNormal` | 1 | Standard notification with sound |
| `PriorityHigh` | 2 | Auto-mentions the user |
| `PriorityCritical` | 3 | Mentions user, sent to all channels, never delayed |

> **Global menu button:** If a notification message has no `Options` block and the host has a menu command configured (typically `/start`), a "🏠 Главное меню" button is automatically appended.

### Single user — plain text

```go
err := ctx.NotifyUser(userID int64, text string, priority int) error
```

### Single user — rich message

```go
err := ctx.NotifyRecipient(userID int64, msg Message, priority int) error
```

### Multiple users — rich message

```go
err := ctx.NotifyUsers(userIDs []int64, msg Message, priority int) error
```

### Teacher notifications

Target a teacher by any of the three university identifiers:

```go
// By teacher_positions.id
err := ctx.NotifyTeacher(teacherPositionID int64, msg Message, priority int) error

// By persons.id
err := ctx.NotifyTeacherPerson(personID int64, msg Message, priority int) error

// By persons.external_id (TSU account ID)
err := ctx.NotifyTeacherExternalID(externalID string, msg Message, priority int) error
```

### Chat notification

Sends to a specific messenger chat (not a user):

```go
err := ctx.NotifyChat(channelType, chatID, text string, priority int) error
```

### Builder: multiple recipients

```go
err := ctx.NotifyRecipients().
    User(userID).
    Users(id1, id2, id3).
    Teacher(teacherPositionID).
    TeacherPerson(personID).
    TeacherExternalID(externalID).
    Message(wasmplugin.NewMessage("Уведомление")).
    Priority(wasmplugin.PriorityHigh).
    Send()
```

### Builder: students by university scope

```go
err := ctx.NotifyStudents().
    Stream(streamID).      // or Faculty / Department / Program / Group / Subgroup
    Message(wasmplugin.NewMessage("Пары завтра отменены")).
    Priority(wasmplugin.PriorityHigh).
    Send()
```

Available scope methods: `Faculty(id)`, `Department(id)`, `Program(id)`, `Stream(id)`, `Group(id)`, `Subgroup(id)`.

---

## Database

Requires `wasmplugin.Database("...").Build()` in `Requirements`.

The plugin gets its own isolated SQL database. Use standard `database/sql`:

```go
import "database/sql"
import wasmplugin "github.com/SuperBotForge/sdk/go-sdk" // driver registers itself via init()

db, err := sql.Open("superbot", "")          // default database
db, err := sql.Open("superbot", "analytics") // named database
```

Declare SQL schema migrations in `Plugin.Migrations`:

```go
Migrations: []wasmplugin.SQLMigration{
    {
        Version:     1,
        Description: "create_queue_entries",
        Up: `CREATE TABLE queue_entries (
            id SERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )`,
        Down: `DROP TABLE queue_entries`,
    },
},
```

Migrations run automatically via goose before `OnConfigure` is called.

---

## File Store

Requires `wasmplugin.File("...").Build()` in `Requirements`.

### Read metadata

```go
ref, err := ctx.FileMeta(fileID string) (*FileRef, error)
```

### Read content

```go
// Read a chunk (up to 1 MB per call)
data, eof, err := ctx.FileRead(fileID string, offset int64, maxBytes int64) ([]byte, bool, error)

// Read entire file (small files only)
data, err := ctx.FileReadAll(fileID string) ([]byte, error)
```

### Get download URL

```go
url, err := ctx.FileURL(fileID string) (string, error)
```

### Store a file

```go
ref, err := ctx.FileStore(name, mimeType, fileType string, data []byte) (*FileRef, error)

// With TTL
ref, err := ctx.FileStoreWithTTL(name, mimeType, fileType string, data []byte, ttl time.Duration) (*FileRef, error)
```

`FileRef` fields:

```go
type FileRef struct {
    ID       string
    Name     string
    MIMEType string
    Size     int64
    FileType string // e.g. "document", "image", "audio"
}
```

Incoming file references from messenger events are in `ctx.Messenger.Files` and can be passed directly to notify/reply message builders.

---

## Inter-plugin RPC

Requires `wasmplugin.PluginDep("target-plugin-id", "...").Build()` in `Requirements`.

### Exposing an RPC method

```go
RPCMethods: []wasmplugin.RPCMethod{
    {
        Name:        "get_status",
        Description: "Returns current queue status",
        Handler: func(ctx *wasmplugin.RPCContext) ([]byte, error) {
            var params struct{ QueueID int64 }
            if err := ctx.Decode(&params); err != nil {
                return nil, err
            }
            result := map[string]interface{}{"length": 42}
            return wasmplugin.MarshalRPC(result)
        },
    },
},
```

### Calling another plugin

```go
// From EventContext — host function (not in SDK yet; use http trigger or events)
```

`RPCContext` fields:

```go
ctx.Caller string  // plugin ID of the caller
ctx.Method string  // method name
ctx.Params []byte  // msgpack-encoded parameters
ctx.Decode(&v)     // decode params into v
ctx.Config("key", "fallback")
ctx.Log("msg")
ctx.LogError("msg")
```

---

## Plugin struct reference

```go
wasmplugin.Plugin{
    ID:      "unique-plugin-id",
    Name:    "Human-readable name",
    Version: "1.2.3",

    // Config schema — admin UI generates a form from this
    Config: wasmplugin.ConfigFields(
        wasmplugin.String("api_key", "API key").Required(),
        wasmplugin.Integer("timeout", "Timeout (seconds)").Default(30).Min(1).Max(300),
        wasmplugin.Bool("verbose", "Verbose logging"),
        wasmplugin.Enum("mode", "Operation mode", "fast", "safe"),
        wasmplugin.StringArray("allowed_hosts", "Allowed hosts"),
    ),

    Requirements: []wasmplugin.Requirement{
        wasmplugin.Database("Store queue entries").Build(),
        wasmplugin.KV("Session state").Build(),
        wasmplugin.NotifyReq("Notify users").Build(),
        wasmplugin.HTTP("Call university API").Build(),
        wasmplugin.UserInfoReq("Read user profiles").Build(),
        wasmplugin.File("Handle attachments").Build(),
        wasmplugin.PluginDep("other-plugin-id", "Get schedule").Build(),
    },

    Migrations: []wasmplugin.SQLMigration{ /* ... */ },

    Triggers: []wasmplugin.Trigger{ /* ... */ },

    OnConfigure:   func(config []byte) error { /* ... */ },
    OnReconfigure: func(old, new []byte) error { /* ... */ },
    OnEvent:       func(ctx *wasmplugin.EventContext) error { /* ... */ },
    Migrate:       func(ctx *wasmplugin.MigrateContext) error { /* ... */ },

    // Return visible command names for a given user (nil = all visible)
    CheckVisibility: func(ctx *wasmplugin.VisibilityContext) []string {
        if isAdmin(ctx.UserID) {
            return nil // all visible
        }
        return []string{"/hello", "/status"}
    },
}
```

---

## Trigger types

### Messenger command

```go
{
    Name: "/start",
    Type: wasmplugin.TriggerMessenger,
    Descriptions: map[string]string{
        "ru": "Начать работу",
        "en": "Get started",
    },
    Handler: handleStart,
}
```

### HTTP endpoint

```go
{
    Name:    "webhook",
    Type:    wasmplugin.TriggerHTTP,
    Path:    "/webhook",
    Methods: []string{"POST"},
    Handler: handleWebhook,
}
```

### Cron schedule

```go
{
    Name:     "daily-digest",
    Type:     wasmplugin.TriggerCron,
    Schedule: "0 8 * * 1-5", // 08:00 on weekdays
    Handler:  handleDailyDigest,
}
```

### Event bus

```go
{
    Name:    "user-linked",
    Type:    wasmplugin.TriggerEvent,
    Topic:   "user.tsu_linked",
    Handler: handleUserLinked,
}
```

---

## Version migration

The `Migrate` callback runs when the host detects a version change during plugin reload. Use it to transform KV data between versions.

```go
Migrate: func(ctx *wasmplugin.MigrateContext) error {
    ctx.OldVersion // e.g. "1.0.0"
    ctx.NewVersion // e.g. "2.0.0"

    // Read, transform, re-write KV data
    keys, _ := ctx.KVList("session:")
    for _, key := range keys {
        val, found, _ := ctx.KVGet(key)
        if !found {
            continue
        }
        // ... transform val ...
        ctx.KVSet(key, newVal)
    }
    return nil
},
```

SQL schema migrations (`Plugin.Migrations`) run automatically before this callback — you don't need to handle them here.
