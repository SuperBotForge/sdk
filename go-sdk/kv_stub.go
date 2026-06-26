//go:build !wasip1

package wasmplugin

import (
	"fmt"
	"time"
)

func (ctx *EventContext) KVGet(key string) (string, bool, error) {
	return "", false, fmt.Errorf("KVGet is only available in WASM")
}

func (ctx *EventContext) KVSet(key, value string) error {
	return fmt.Errorf("KVSet is only available in WASM")
}

func (ctx *EventContext) KVSetWithTTL(key, value string, ttl time.Duration) error {
	return fmt.Errorf("KVSetWithTTL is only available in WASM")
}

func (ctx *EventContext) KVDelete(key string) error {
	return fmt.Errorf("KVDelete is only available in WASM")
}

func (ctx *EventContext) KVList(prefix string) ([]string, error) {
	return nil, fmt.Errorf("KVList is only available in WASM")
}

func (ctx *MigrateContext) KVGet(key string) (string, bool, error) {
	return "", false, fmt.Errorf("KVGet is only available in WASM")
}

func (ctx *MigrateContext) KVSet(key, value string) error {
	return fmt.Errorf("KVSet is only available in WASM")
}

func (ctx *MigrateContext) KVDelete(key string) error {
	return fmt.Errorf("KVDelete is only available in WASM")
}

func (ctx *MigrateContext) KVList(prefix string) ([]string, error) {
	return nil, fmt.Errorf("KVList is only available in WASM")
}
