package chat

// Storage common interface
type Storage interface {
	Open() error
	Close()

	PushMessage(chatID int64, msg *Message) error
	AddRoom(chatID int64, title string) error

	LoadRooms() []*Room
	LoadMessages() map[string][]*Message // [chatID]messages

	DeleteMessage(chatID int64, msgID int64) error
	DeleteChat(chatID int64) error

	LoadProfile() *Profile
	SetUserName(userName string) error
}
