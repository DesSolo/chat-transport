package vk

import (
	"bytes"
	"chat-transport/internal/entities"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

// Client ...
type Client struct {
	name           string
	accessToken    string
	chatID         int
	strChatID      string
	ignoreAccounts []int
	template       *template.Template
	lastMessageID  int
	client         *http.Client
}

// NewClient ...
func NewClient(name, accessToken string, chatID int, ignoreAccounts []int, t *template.Template) *Client {
	return &Client{
		name:           name,
		accessToken:    accessToken,
		chatID:         chatID,
		strChatID:      strconv.Itoa(chatID),
		ignoreAccounts: ignoreAccounts,
		template:       t,
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
	if c.accessToken == "" {
		return errors.New("token value not valid")
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
		"extended": 1,
		"fields":   "first_name,last_name",
		"peer_id":  c.chatID,
	}
	data, err := c.callMethod("messages.getHistory", p)
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

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	profiles := make(map[int]string)
	for _, pr := range resp.Response.Profiles {
		profiles[pr.ID] = fmt.Sprintf("%s %s", pr.FirstName, pr.LastName)
	}

	var messages []*entities.Message

	for i := len(resp.Response.Items) - 1; i >= 0; i-- {
		item := resp.Response.Items[i]

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
	var msg bytes.Buffer

	if err := c.template.Execute(&msg, m); err != nil {
		return err
	}

	p := RequestParams{
		"peer_id":   c.chatID,
		"message":   msg.String(),
		"random_id": 0,
	}

	_, err := c.callMethod("messages.send", p)
	if err != nil {
		return err
	}

	return nil
}
