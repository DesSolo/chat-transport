package telegram

import (
	"chat-transport/internal/entities"
	"errors"
	"net/http"
	"strconv"
)

// Telegram ...
type Telegram struct {
	name           string
	token          string
	chatID         string
	ignoreAccounts []string
	lastUpdateID   int
	client         *http.Client
}

// NewTelegram ...
func NewTelegram(name, token, chatID string, ignoreAccounts []string) *Telegram {
	return &Telegram{
		name:           name,
		token:          token,
		chatID:         chatID,
		ignoreAccounts: ignoreAccounts,
		client:         &http.Client{},
	}
}

// GetName ...
func (t *Telegram) GetName() string {
	return t.name
}

// GetChatID ...
func (t *Telegram) GetChatID() string {
	return t.chatID
}

// Validate ...
func (t *Telegram) Validate() error {
	if t.token == "" {
		return errors.New("token value not valid")
	}

	if t.chatID == "" {
		return errors.New("chat value not valid")
	}

	return nil
}

func (t *Telegram) isIgnoreUser(username string) bool {
	for _, ia := range t.ignoreAccounts {
		if ia == username {
			return true
		}
	}

	return false
}

// GetNewMessages ...
func (t *Telegram) GetNewMessages() ([]*entities.Message, error) {
	updates, err := t.getUpdates(t.lastUpdateID)
	if err != nil {
		return nil, err
	}

	var messages []*entities.Message
	for _, upd := range updates {

		t.lastUpdateID = upd.UpdateID

		if upd.Message.Text == "" {
			continue
		}

		chatID := strconv.Itoa(upd.Message.Chat.ID)

		if chatID != t.chatID {
			continue
		}

		if t.isIgnoreUser(upd.Message.From.Username) {
			continue
		}

		msg := entities.Message{
			Author: entities.Author{
				Username: upd.Message.From.Username,
			},
			Text: upd.Message.Text,
			Chat: entities.Chat{
				ID: chatID,
			},
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

// SendMessage ...
func (t *Telegram) SendMessage(m *entities.Message) error {
	return t.sendMessage(t.chatID, m.Text)
}
