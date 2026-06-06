//go:build !wasip1

package wasmplugin

import "fmt"

func GetUserInfo(userID int64) (*UserInfo, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("user_id must be greater than zero")
	}
	return nil, fmt.Errorf("GetUserInfo is only available in WASM")
}

func (ctx *EventContext) GetUserInfo(userID int64) (*UserInfo, error) {
	return GetUserInfo(userID)
}
