package protocol_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type schemaHeader struct {
	ID string `json:"$id"`
}

func TestSchemasCompile(t *testing.T) {
	compileSchemas(t)
}

func TestValidFixturesConformToSchemas(t *testing.T) {
	schemas := compileSchemas(t)

	tests := []struct {
		name       string
		schemaFile string
		fixture    string
	}{
		{"plugin meta", "plugin-meta.schema.json", "testdata/valid/plugin-meta.json"},
		{"messenger event request", "event-request.schema.json", "testdata/valid/event-request-messenger.json"},
		{"http event request", "event-request.schema.json", "testdata/valid/event-request-http.json"},
		{"event response", "event-response.schema.json", "testdata/valid/event-response.json"},
		{"rpc request", "rpc-request.schema.json", "testdata/valid/rpc-request.json"},
		{"rpc response", "rpc-response.schema.json", "testdata/valid/rpc-response.json"},
		{"step callback request", "step-callback-request.schema.json", "testdata/valid/step-callback-request.json"},
		{"step callback response", "step-callback-response.schema.json", "testdata/valid/step-callback-response.json"},
		{"migrate request", "migrate-request.schema.json", "testdata/valid/migrate-request.json"},
		{"migrate response", "migrate-response.schema.json", "testdata/valid/migrate-response.json"},
		{"reconfigure request", "reconfigure-request.schema.json", "testdata/valid/reconfigure-request.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateFixture(t, schemas[tt.schemaFile], tt.fixture)
		})
	}
}

func TestInvalidFixturesDoNotConformToSchemas(t *testing.T) {
	schemas := compileSchemas(t)

	tests := []struct {
		name       string
		schemaFile string
		fixture    string
	}{
		{"plugin meta wrong version", "plugin-meta.schema.json", "testdata/invalid/plugin-meta-wrong-version.json"},
		{"unknown event trigger", "event-request.schema.json", "testdata/invalid/event-request-unknown-trigger.json"},
		{"rpc request missing method", "rpc-request.schema.json", "testdata/invalid/rpc-request-missing-method.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := readJSONDocument(t, tt.fixture)
			if err := schemas[tt.schemaFile].Validate(doc); err == nil {
				t.Fatalf("%s unexpectedly conformed to %s", tt.fixture, tt.schemaFile)
			}
		})
	}
}

func compileSchemas(t *testing.T) map[string]*jsonschema.Schema {
	t.Helper()

	matches, err := filepath.Glob("schemas/*.schema.json")
	if err != nil {
		t.Fatalf("glob schemas: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("no schema files found")
	}

	compiler := jsonschema.NewCompiler()
	idsByFile := make(map[string]string, len(matches))

	for _, path := range matches {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read schema %s: %v", path, err)
		}

		var header schemaHeader
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

func validateFixture(t *testing.T, schema *jsonschema.Schema, fixture string) {
	t.Helper()
	if schema == nil {
		t.Fatalf("schema for %s is nil", fixture)
	}
	doc := readJSONDocument(t, fixture)
	if err := schema.Validate(doc); err != nil {
		t.Fatalf("%s does not conform: %v", fixture, err)
	}
}

func readJSONDocument(t *testing.T, path string) any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("unmarshal fixture %s: %v", path, err)
	}
	return doc
}
