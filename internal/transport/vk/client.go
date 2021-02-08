package vk

import (
	"chat-transport/internal/entities"
	"errors"
	"fmt"
	"net/http"
)

// Client ...
type Client struct {
	name           string
	accessToken    string
	chatID         string
	ignoreAccounts []string
	lastUpdateID   int
	client         *http.Client
}

// NewClient ...
func NewClient(name, accessToken, chatID string, ignoreAccounts []string) *Client {
	return &Client{
		name:           name,
		accessToken:    accessToken,
		chatID:         chatID,
		ignoreAccounts: ignoreAccounts,
		client:         &http.Client{},
	}
}

// GetName ...
func (c *Client) GetName() string {
	return c.name
}

// GetChatID ...
func (c *Client) GetChatID() string {
	return c.chatID
}

// Validate ...
func (c *Client) Validate() error {
	if c.accessToken == "" {
		return errors.New("access_token value not valid")
	}

	if c.chatID == "" {
		return errors.New("chat value not valid")
	}

	return nil
}

// GetNewMessages ...
func (c *Client) GetNewMessages() ([]*entities.Message, error) {
	// https://vk.com/dev/messages.get
	p := RequestParams{
		"offset":          1,
		"time_offset":     "",
		"last_message_id": 0,
	}
	data, err := c.callMethod("messages.get", p)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s", data)

	return nil, nil
}

// SendMessage ...
func (c *Client) SendMessage(m *entities.Message) error {
	// https://vk.com/dev/messages.send
	p := RequestParams{
		"user_id": c.chatID,
		"peer_id": c.chatID,
		"chat_id": c.chatID,
		"message": m.Text,
	}
	data, err := c.callMethod("messages.send", p)
	if err != nil {
		return err
	}

	fmt.Printf("%s", data)

	return nil
}
