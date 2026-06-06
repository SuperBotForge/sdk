//go:build wasip1

package wasmplugin

import "fmt"

//go:wasmimport env user_info
func _userInfo(offset, length uint32) uint64

type userInfoReq struct {
	UserID int64 `msgpack:"user_id"`
}

// GetUserInfo fetches information about a user by their global user ID.
func GetUserInfo(userID int64) (*UserInfo, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("user_id must be greater than zero")
	}

	var info UserInfo
	if err := callHostWithResult(_userInfo, userInfoReq{UserID: userID}, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetUserInfo fetches user information using the current event context.
func (ctx *EventContext) GetUserInfo(userID int64) (*UserInfo, error) {
	return GetUserInfo(userID)
}
