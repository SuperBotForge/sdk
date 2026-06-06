package wasmplugin

// UserInfo contains basic information about a user.
// The struct is intentionally minimal and will grow over time.
type UserInfo struct {
	ID       int64  `msgpack:"id" json:"id"`
	FullName string `msgpack:"full_name,omitempty" json:"full_name,omitempty"`
}
