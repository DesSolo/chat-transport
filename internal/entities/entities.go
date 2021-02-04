package entities

// Author ..
type Author struct {
	Username string
}

// Message ..
type Message struct {
	Author Author
	Text   string
}

// Transport ...
type Transport interface {
	GetName() string
	Validate() error
	GetNewMessages() ([]*Message, error)
	SendMessage(*Message) error
}
