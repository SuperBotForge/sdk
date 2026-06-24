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

func GetUsersInfo(userIDs []int64) ([]UserInfoFull, error) {
	return nil, fmt.Errorf("GetUsersInfo is only available in WASM")
}

func (ctx *EventContext) GetUsersInfo(userIDs []int64) ([]UserInfoFull, error) {
	return GetUsersInfo(userIDs)
}

func ListUsers(page, pageSize int) ([]UserInfoFull, int, error) {
	return nil, 0, fmt.Errorf("ListUsers is only available in WASM")
}

func (ctx *EventContext) ListUsers(page, pageSize int) ([]UserInfoFull, int, error) {
	return ListUsers(page, pageSize)
}
