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
	authBasic      string
	cookieRemixsID string
	chatID         string
	ignoreAccounts []string
	// lastUpdateID   int
	client *http.Client
}

// NewClient ...
func NewClient(name, authBasic, cookieRemixsID, chatID string, ignoreAccounts []string) *Client {
	return &Client{
		name:           name,
		authBasic:      authBasic,
		cookieRemixsID: cookieRemixsID,
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
	if c.authBasic == "" {
		return errors.New("auth_basic value not valid")
	}

	if c.cookieRemixsID == "" {
		return errors.New("cookie_remixsid not valid")
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
		"filters": 8,
		"offset":  0,
	}
	data, err := c.callMethod("messages.get", "1612984943:09d70e65be95b44b7a", p)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Response struct {
			Count int      `json:"count"`
			Items []string `json:"items"`
		} `json:"response"`
	}

	if err := unmarshal(data, &resp); err != nil {
		return nil, err
	}

	fmt.Printf("total %d", resp.Response.Count)

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
	data, err := c.callMethod("messages.send", "", p)
	if err != nil {
		return err
	}

	fmt.Printf("%s", data)

	return nil
}
