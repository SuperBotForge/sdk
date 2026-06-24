package wasmplugin

// UserInfo contains basic information about a user.
type UserInfo struct {
	ID         int64  `msgpack:"id" json:"id"`
	FullName   string `msgpack:"full_name,omitempty" json:"full_name,omitempty"`
	ExternalID string `msgpack:"external_id,omitempty" json:"external_id,omitempty"`
	IsTeacher  bool   `msgpack:"is_teacher,omitempty" json:"is_teacher,omitempty"`
}

// UserPosition holds one position (student/teacher) for a user.
type UserPosition struct {
	PositionType    string `msgpack:"position_type" json:"position_type"`
	Status          string `msgpack:"status,omitempty" json:"status,omitempty"`
	NationalityType string `msgpack:"nationality_type,omitempty" json:"nationality_type,omitempty"`
	FundingType     string `msgpack:"funding_type,omitempty" json:"funding_type,omitempty"`
	EducationForm   string `msgpack:"education_form,omitempty" json:"education_form,omitempty"`
	GroupCode       string `msgpack:"group_code,omitempty" json:"group_code,omitempty"`
	GroupName       string `msgpack:"group_name,omitempty" json:"group_name,omitempty"`
	ProgramName     string `msgpack:"program_name,omitempty" json:"program_name,omitempty"`
	StreamName      string `msgpack:"stream_name,omitempty" json:"stream_name,omitempty"`
}

// UserInfoFull extends UserInfo with university positions.
type UserInfoFull struct {
	UserInfo
	Positions []UserPosition `msgpack:"positions,omitempty" json:"positions,omitempty"`
}
