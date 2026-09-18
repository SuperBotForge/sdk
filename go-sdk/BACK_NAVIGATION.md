# Dialog back navigation

Messenger triggers can opt into host-managed back navigation:

```go
Trigger{
    Name: "tasks",
    Type: TriggerMessenger,
    AllowBack: true,
    Nodes: []Node{/* steps */},
    Handler: handleTasks,
}
```

The metadata field is `allow_back` (optional boolean, false by default).
Install a Core build with navigation support before enabling it in a plugin.
The protocol remains v4; its metadata schema now includes this optional field.

Core adds a localized Back option to each unfinished step, including text and
file prompts. It restores the answers and pagination from before the previous
accepted answer. Hidden steps are not recorded. Pagination and validation
failures do not add history. Changing a branch therefore discards the abandoned
branch's answers. Previously entered values on the revisited step must be entered
again.

Back at the first step cancels the command and opens the plugin's command menu.
The plugin's event handler is not called on cancellation. Completed actions are
not undone. Commands with no steps run immediately and have no back button.

The host persists history with the dialog. A token prevents stale Back callbacks
from changing a newer step or another dialog. Do not use the reserved callback
prefix `__dialog_back:` as an option value. Callback functions should remain
read-only; navigation does not undo external side effects performed by callbacks.

This feature does not change file transport: the existing host sends attachments
from the completing input to the plugin. Place file uploads on the final step
when the handler needs the uploaded file. Navigation before that upload is supported.
