//go:build !wasip1

package wasmplugin

import "fmt"

const (
	PriorityLow      = 0
	PriorityNormal   = 1
	PriorityHigh     = 2
	PriorityCritical = 3
)

func (ctx *EventContext) NotifyUser(userID int64, text string, priority int) error {
	return fmt.Errorf("NotifyUser is only available in WASM")
}

func (ctx *EventContext) NotifyRecipient(userID int64, msg Message, priority int) error {
	return fmt.Errorf("NotifyRecipient is only available in WASM")
}

func (ctx *EventContext) NotifyTeacher(teacherPositionID int64, msg Message, priority int) error {
	return fmt.Errorf("NotifyTeacher is only available in WASM")
}

func (ctx *EventContext) NotifyTeacherPerson(personID int64, msg Message, priority int) error {
	return fmt.Errorf("NotifyTeacherPerson is only available in WASM")
}

func (ctx *EventContext) NotifyTeacherExternalID(externalID string, msg Message, priority int) error {
	return fmt.Errorf("NotifyTeacherExternalID is only available in WASM")
}

func (ctx *EventContext) NotifyChat(channelType, chatID, text string, priority int) error {
	return fmt.Errorf("NotifyChat is only available in WASM")
}

func (ctx *EventContext) NotifyUsers(userIDs []int64, msg Message, priority int) error {
	return fmt.Errorf("NotifyUsers is only available in WASM")
}

func (ctx *EventContext) NotifyRecipients() *RecipientNotifyBuilder {
	return &RecipientNotifyBuilder{}
}

func (ctx *EventContext) NotifyStudents() *StudentNotifyBuilder {
	return &StudentNotifyBuilder{}
}

type RecipientNotifyBuilder struct{}

func (b *RecipientNotifyBuilder) User(userID int64) *RecipientNotifyBuilder      { return b }
func (b *RecipientNotifyBuilder) Users(userIDs ...int64) *RecipientNotifyBuilder  { return b }
func (b *RecipientNotifyBuilder) Teacher(id int64) *RecipientNotifyBuilder        { return b }
func (b *RecipientNotifyBuilder) TeacherPerson(id int64) *RecipientNotifyBuilder  { return b }
func (b *RecipientNotifyBuilder) TeacherExternalID(id string) *RecipientNotifyBuilder { return b }
func (b *RecipientNotifyBuilder) Message(msg Message) *RecipientNotifyBuilder     { return b }
func (b *RecipientNotifyBuilder) Priority(p int) *RecipientNotifyBuilder          { return b }
func (b *RecipientNotifyBuilder) Send() error {
	return fmt.Errorf("NotifyRecipients is only available in WASM")
}

type StudentNotifyBuilder struct{}

func (b *StudentNotifyBuilder) Faculty(id int64) *StudentNotifyBuilder    { return b }
func (b *StudentNotifyBuilder) Department(id int64) *StudentNotifyBuilder { return b }
func (b *StudentNotifyBuilder) Program(id int64) *StudentNotifyBuilder    { return b }
func (b *StudentNotifyBuilder) Stream(id int64) *StudentNotifyBuilder     { return b }
func (b *StudentNotifyBuilder) Group(id int64) *StudentNotifyBuilder      { return b }
func (b *StudentNotifyBuilder) Subgroup(id int64) *StudentNotifyBuilder   { return b }
func (b *StudentNotifyBuilder) Message(msg Message) *StudentNotifyBuilder { return b }
func (b *StudentNotifyBuilder) Priority(p int) *StudentNotifyBuilder      { return b }
func (b *StudentNotifyBuilder) Send() error {
	return fmt.Errorf("NotifyStudents is only available in WASM")
}
