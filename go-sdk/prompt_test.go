package wasmplugin

import (
	"encoding/json"
	"testing"
)

func TestPromptCallbackContract(t *testing.T) {
	p := contractTestPlugin()
	step := NewStep("choice").LocalizedOptions(map[string]string{"en": "Fallback"}, Opt("Yes", "yes")).PromptsFn(func(c *CallbackContext) map[string]string {
		return map[string]string{"ru": "ru:" + c.Params["semester"], "en": "en:" + c.Params["semester"]}
	})
	step.Text("fallback", StylePlain).PromptsFn(func(c *CallbackContext) map[string]string { return map[string]string{"ru": ""} })
	p.Triggers[0].Nodes = []Node{step}
	reg := callbackMap{}
	nd := step.toNodeDef(p.Triggers[0].Name, reg)
	if nd.Blocks[0].PromptsFn == nd.Blocks[1].PromptsFn {
		t.Fatal("callback name collision")
	}
	meta := runPluginAction(t, p, "meta", nil, nil)
	validateProtocolJSON(t, loadProtocolSchemas(t)["plugin-meta.schema.json"], meta)
	for i, want := range []string{"ru:5", ""} {
		request, _ := json.Marshal(stepCallbackRequest{Callback: nd.Blocks[i].PromptsFn, Locale: "ru", Params: map[string]string{"semester": "5"}})
		raw := runPluginAction(t, p, "step_callback", request, nil)
		validateProtocolJSON(t, loadProtocolSchemas(t)["step-callback-response.schema.json"], raw)
		var resp stepCallbackResponse
		if err := json.Unmarshal(raw, &resp); err != nil {
			t.Fatal(err)
		}
		if resp.Prompts == nil || resp.Prompts["ru"] != want {
			t.Fatalf("response: %s", raw)
		}
	}
}
