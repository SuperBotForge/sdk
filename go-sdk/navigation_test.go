package wasmplugin

import (
	"encoding/json"
	"testing"
)

func TestAllowBackMetadata(t *testing.T) {
	p := contractTestPlugin()
	p.Triggers[0].AllowBack = true
	raw := runPluginAction(t, p, "meta", nil, nil)
	validateProtocolJSON(t, loadProtocolSchemas(t)["plugin-meta.schema.json"], raw)
	var meta pluginMeta
	if err := json.Unmarshal(raw, &meta); err != nil {
		t.Fatal(err)
	}
	if !meta.Triggers[0].AllowBack {
		t.Fatal("AllowBack missing from metadata")
	}
}
