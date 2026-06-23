package wasmplugin

// UserInfo contains basic information about a user.
type UserInfo struct {
	ID         int64  `msgpack:"id" json:"id"`
	FullName   string `msgpack:"full_name,omitempty" json:"full_name,omitempty"`
	ExternalID string `msgpack:"external_id,omitempty" json:"external_id,omitempty"`
	IsTeacher  bool   `msgpack:"is_teacher,omitempty" json:"is_teacher,omitempty"`
}
