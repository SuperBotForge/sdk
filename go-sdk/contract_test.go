package wasmplugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type protocolSchemaHeader struct {
	ID string `json:"$id"`
}

func TestGoSDKLifecycleResponsesConformToProtocol(t *testing.T) {
	schemas := loadProtocolSchemas(t)
	plugin := contractTestPlugin()

	tests := []struct {
		name           string
		action         string
		requestSchema  string
		request        string
		config         string
		responseSchema string
		plugin         Plugin
	}{
		{
			name:           "meta",
			action:         "meta",
			responseSchema: "plugin-meta.schema.json",
			plugin:         plugin,
		},
		{
			name:          "handle_event",
			action:        "handle_event",
			requestSchema: "event-request.schema.json",
			request: `{
				"id": "evt-sdk-1",
				"trigger_type": "messenger",
				"trigger_name": "hello",
				"plugin_id": "contract",
				"timestamp": 1710000000,
				"data": {
					"user_id": 42,
					"channel_type": "telegram",
					"chat_id": "chat-1",
					"command_name": "hello",
					"params": {"choice": "yes"},
					"locale": "en"
				}
			}`,
			config:         `{"greeting":"Hi"}`,
			responseSchema: "event-response.schema.json",
			plugin:         plugin,
		},
		{
			name:          "handle_rpc",
			action:        "handle_rpc",
			requestSchema: "rpc-request.schema.json",
			request: `{
				"caller": "tester",
				"method": "echo",
				"params": "gqJpZAQ="
			}`,
			config:         `{"greeting":"Hi"}`,
			responseSchema: "rpc-response.schema.json",
			plugin:         plugin,
		},
		{
			name:          "step_callback",
			action:        "step_callback",
			requestSchema: "step-callback-request.schema.json",
			request: `{
				"callback": "hello:options:choice",
				"user_id": 42,
				"locale": "en",
				"params": {"choice": "yes"},
				"page": 0,
				"input": ""
			}`,
			config:         `{"greeting":"Hi"}`,
			responseSchema: "step-callback-response.schema.json",
			plugin:         plugin,
		},
		{
			name:          "migrate",
			action:        "migrate",
			requestSchema: "migrate-request.schema.json",
			request: `{
				"old_version": "1.0.0",
				"new_version": "1.1.0"
			}`,
			responseSchema: "migrate-response.schema.json",
			plugin:         plugin,
		},
		{
			name:          "reconfigure error",
			action:        "reconfigure",
			requestSchema: "reconfigure-request.schema.json",
			request: `{
				"previous_config": {"greeting": "Hello"},
				"config": {"greeting": "Hi"}
			}`,
			responseSchema: "event-response.schema.json",
			plugin: Plugin{
				ID:      "contract",
				Name:    "Contract Plugin",
				Version: "1.0.0",
				OnReconfigure: func(previousConfig, config []byte) error {
					return errors.New("reconfigure failed")
				},
			},
		},
		{
			name:           "configure error",
			action:         "configure",
			requestSchema:  "",
			request:        `{"greeting":"Hi"}`,
			responseSchema: "event-response.schema.json",
			plugin: Plugin{
				ID:      "contract",
				Name:    "Contract Plugin",
				Version: "1.0.0",
				OnConfigure: func(config []byte) error {
					return errors.New("configure failed")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requestSchema != "" {
				validateProtocolJSON(t, schemas[tt.requestSchema], []byte(tt.request))
			}

			stdout := runPluginAction(t, tt.plugin, tt.action, []byte(tt.request), []byte(tt.config))
			if len(stdout) == 0 {
				t.Fatalf("%s wrote empty stdout, want %s", tt.action, tt.responseSchema)
			}
			validateProtocolJSON(t, schemas[tt.responseSchema], stdout)
		})
	}
}

func TestGoSDKSuccessfulConfigureWritesNoProtocolBody(t *testing.T) {
	stdout := runPluginAction(t, Plugin{
		ID:      "contract",
		Name:    "Contract Plugin",
		Version: "1.0.0",
		OnConfigure: func(config []byte) error {
			return nil
		},
	}, "configure", []byte(`{"greeting":"Hi"}`), nil)

	if len(stdout) != 0 {
		t.Fatalf("configure stdout = %s, want empty success body", stdout)
	}
}

func contractTestPlugin() Plugin {
	return Plugin{
		ID:      "contract",
		Name:    "Contract Plugin",
		Version: "1.0.0",
		Config: ConfigFields(
			String("greeting", "Greeting").Default("Hello"),
		),
		Requirements: []Requirement{
			Database("Store contract test data").Build(),
			HTTP("Fetch test data").WithConfig(HTTPPolicyConfig()).Build(),
			KV("Store lightweight values").Build(),
			NotifyReq("Send notifications").Build(),
			EventsReq("Publish events").Build(),
			PluginDep("core", "Call core plugin").Build(),
			File("Read and write files").Build(),
		},
		RPCMethods: []RPCMethod{
			{
				Name:        "echo",
				Description: "Echo test payload",
				Handler: func(ctx *RPCContext) ([]byte, error) {
					ctx.Log("rpc handled")
					return MarshalRPC(map[string]bool{"ok": true})
				},
			},
		},
		Triggers: []Trigger{
			{
				Name: "hello",
				Type: TriggerMessenger,
				Descriptions: map[string]string{
					"en": "Hello",
					"ru": "Privet",
				},
				Nodes: []Node{
					NewStep("choice").
						LocalizedText(map[string]string{"en": "Choose"}, StyleHeader).
						LocalizedDynamicOptions(map[string]string{"en": "Choice"}, func(ctx *CallbackContext) []Option {
							return []Option{
								{Label: "Yes", Labels: map[string]string{"en": "Yes"}, Value: "yes"},
								{Label: "No", Labels: map[string]string{"en": "No"}, Value: "no"},
							}
						}).
						ValidateFunc(func(ctx *CallbackContext) bool {
							return ctx.Input == "yes" || ctx.Input == "no"
						}).
						VisibleWhen(ParamNeq("hidden", "true")),
				},
				Handler: func(ctx *EventContext) error {
					ctx.Log("event handled")
					ctx.Reply(NewMessage(ctx.Config("greeting", "Hello")))
					return nil
				},
			},
			{
				Name:    "webhook",
				Type:    TriggerHTTP,
				Path:    "/webhook",
				Methods: []string{"POST"},
			},
			{
				Name:     "daily",
				Type:     TriggerCron,
				Schedule: "0 8 * * *",
			},
			{
				Name:  "sync",
				Type:  TriggerEvent,
				Topic: "contract.sync",
			},
		},
		Migrate: func(ctx *MigrateContext) error {
			return nil
		},
		Migrations: []SQLMigration{
			{
				Version:     1,
				Description: "init",
				Up:          "CREATE TABLE contract_test(id text);",
				Down:        "DROP TABLE contract_test;",
			},
		},
	}
}

func runPluginAction(t *testing.T, plugin Plugin, action string, input, config []byte) []byte {
	t.Helper()

	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdin pipe: %v", err)
	}
	if _, err := stdinWriter.Write(input); err != nil {
		t.Fatalf("write stdin: %v", err)
	}
	if err := stdinWriter.Close(); err != nil {
		t.Fatalf("close stdin writer: %v", err)
	}

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}

	oldStdin := os.Stdin
	oldStdout := os.Stdout
	os.Stdin = stdinReader
	os.Stdout = stdoutWriter
	t.Cleanup(func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
		_ = stdinReader.Close()
		_ = stdoutReader.Close()
		_ = stdoutWriter.Close()
	})

	t.Setenv("PLUGIN_ACTION", action)
	if len(config) > 0 {
		t.Setenv("PLUGIN_CONFIG", string(config))
	} else {
		t.Setenv("PLUGIN_CONFIG", "")
	}

	Run(plugin)

	os.Stdin = oldStdin
	os.Stdout = oldStdout
	if err := stdoutWriter.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	stdout, err := io.ReadAll(stdoutReader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	return stdout
}

func loadProtocolSchemas(t *testing.T) map[string]*jsonschema.Schema {
	t.Helper()

	schemasDir := findProtocolSchemasDir(t)
	matches, err := filepath.Glob(filepath.Join(schemasDir, "*.schema.json"))
	if err != nil {
		t.Fatalf("glob schemas: %v", err)
	}
	if len(matches) == 0 {
		t.Fatalf("no schema files found in %s", schemasDir)
	}

	compiler := jsonschema.NewCompiler()
	idsByFile := make(map[string]string, len(matches))

	for _, path := range matches {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read schema %s: %v", path, err)
		}

		var header protocolSchemaHeader
		if err := json.Unmarshal(raw, &header); err != nil {
			t.Fatalf("parse schema header %s: %v", path, err)
		}
		if header.ID == "" {
			t.Fatalf("schema %s has no $id", path)
		}

		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
		if err != nil {
			t.Fatalf("unmarshal schema %s: %v", path, err)
		}
		if err := compiler.AddResource(header.ID, doc); err != nil {
			t.Fatalf("add schema resource %s: %v", path, err)
		}
		idsByFile[filepath.Base(path)] = header.ID
	}

	compiled := make(map[string]*jsonschema.Schema, len(idsByFile))
	for file, id := range idsByFile {
		schema, err := compiler.Compile(id)
		if err != nil {
			t.Fatalf("compile schema %s: %v", file, err)
		}
		compiled[file] = schema
	}
	return compiled
}

func findProtocolSchemasDir(t *testing.T) string {
	t.Helper()

	if fromEnv := os.Getenv("SUPERBOTGO_PROTOCOL_SCHEMAS"); fromEnv != "" {
		if dirExists(fromEnv) {
			return fromEnv
		}
		t.Fatalf("SUPERBOTGO_PROTOCOL_SCHEMAS=%q does not exist", fromEnv)
	}

	candidates := []string{
		filepath.Join("..", "protocol", "v4", "schemas"),
		filepath.Join("..", "..", "sdk", "protocol", "v4", "schemas"),
	}
	for _, candidate := range candidates {
		if dirExists(candidate) {
			return candidate
		}
	}

	t.Skip("protocol schemas not found; run from the SuperBotGo monorepo or set SUPERBOTGO_PROTOCOL_SCHEMAS")
	return ""
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func validateProtocolJSON(t *testing.T, schema *jsonschema.Schema, data []byte) {
	t.Helper()
	if schema == nil {
		t.Fatal("schema is nil")
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("unmarshal protocol JSON %s: %v", data, err)
	}
	if err := schema.Validate(doc); err != nil {
		t.Fatalf("protocol JSON does not conform:\n%s\nerror: %v", data, err)
	}
}
