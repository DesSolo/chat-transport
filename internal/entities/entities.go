package entities

// Author ..
type Author struct {
	Username string
}

// Chat ...
type Chat struct {
	ID string
}

// Message ..
type Message struct {
	Chat   Chat
	Author Author
	Text   string
}

// Transport ...
type Transport interface {
	GetName() string
	GetChatID() string
	Validate() error
	GetNewMessages() ([]*Message, error)
	SendMessage(*Message) error
}
