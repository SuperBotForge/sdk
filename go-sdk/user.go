//go:build wasip1

package wasmplugin

import "fmt"

//go:wasmimport env user_info
func _userInfo(offset, length uint32) uint64

//go:wasmimport env users_info
func _usersInfo(offset, length uint32) uint64

type userInfoReq struct {
	UserID int64 `msgpack:"user_id"`
}

type usersInfoReq struct {
	UserIDs []int64 `msgpack:"user_ids"`
}

type usersInfoResult struct {
	Users []UserInfoFull `msgpack:"users"`
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

// GetUsersInfo fetches information for multiple users at once, including their university positions.
func GetUsersInfo(userIDs []int64) ([]UserInfoFull, error) {
	var res usersInfoResult
	if err := callHostWithResult(_usersInfo, usersInfoReq{UserIDs: userIDs}, &res); err != nil {
		return nil, err
	}
	return res.Users, nil
}

// GetUsersInfo fetches information for multiple users using the current event context.
func (ctx *EventContext) GetUsersInfo(userIDs []int64) ([]UserInfoFull, error) {
	return GetUsersInfo(userIDs)
}
