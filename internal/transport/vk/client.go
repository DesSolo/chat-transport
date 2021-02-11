package vk

import (
	"chat-transport/internal/entities"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Client ...
type Client struct {
	name           string
	authBasic      string
	cookieRemixsID string
	chatID         int
	strChatID      string
	ignoreAccounts []int
	lastMessageID  int
	client         *http.Client
}

// NewClient ...
func NewClient(name, authBasic, cookieRemixsID string, chatID int, ignoreAccounts []int) *Client {
	return &Client{
		name:           name,
		authBasic:      authBasic,
		cookieRemixsID: cookieRemixsID,
		chatID:         chatID,
		strChatID:      strconv.Itoa(chatID),
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
	return c.strChatID
}

// Validate ...
func (c *Client) Validate() error {
	if c.authBasic == "" {
		return errors.New("auth_basic value not valid")
	}

	if c.cookieRemixsID == "" {
		return errors.New("cookie_remixsid not valid")
	}

	if c.chatID <= 0 {
		return errors.New("chat value not valid")
	}

	return nil
}

func (c *Client) isIgnore(userID int) bool {
	for _, uid := range c.ignoreAccounts {
		if uid == userID {
			return true
		}
	}

	return false
}

// GetNewMessages ...
func (c *Client) GetNewMessages() ([]*entities.Message, error) {
	// https://vk.com/dev/messages.getHistory
	p := RequestParams{
		"count":    100,
		"extended": 1,
		"fields":   "first_name,last_name",
		"user_id":  c.chatID,
	}
	data, err := c.callMethod("messages.getHistory", "1613066804:bd9ed3fd78469fd161", p)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Response struct {
			Count int `json:"count"`
			Items []struct {
				FromID int    `json:"from_id"`
				ID     int    `json:"id"`
				Text   string `json:"text"`
			} `json:"items"`
			Profiles []struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"profiles"`
		} `json:"response"`
	}

	if err := unmarshal(data, &resp); err != nil {
		return nil, err
	}

	profiles := make(map[int]string)
	for _, pr := range resp.Response.Profiles {
		profiles[pr.ID] = fmt.Sprintf("%s %s", pr.FirstName, pr.LastName)
	}

	var messages []*entities.Message

	for i := len(resp.Response.Items) - 1; i >= 0; i-- {
		item := resp.Response.Items[i]
		if item.FromID != c.chatID {
			continue
		}

		if c.isIgnore(item.FromID) {
			continue
		}

		if c.lastMessageID < item.ID {
			c.lastMessageID = item.ID
			if item.Text == "" {
				continue
			}

			username, ok := profiles[item.FromID]
			if !ok {
				continue
			}

			msg := &entities.Message{
				Chat: entities.Chat{
					ID: c.strChatID,
				},
				Author: entities.Author{
					Username: username,
				},
				Text: item.Text,
			}
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

// SendMessage ...
func (c *Client) SendMessage(m *entities.Message) error {
	// https://vk.com/dev/messages.send
	p := RequestParams{
		"chat_id": c.chatID,
		"message": m.Text,
	}
	data, err := c.callMethod("messages.send", "1612992731:1595b2893c2f8feb1c", p)
	if err != nil {
		return err
	}

	fmt.Printf("%s", data)

	return nil
}
