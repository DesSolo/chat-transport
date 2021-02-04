package transport

import (
	"bytes"
	"chat-transport/internal/entities"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Update ...
type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Date int `json:"date"`
		Chat struct {
			ID    int    `json:"id"`
			Type  string `json:"type"`
			Title string `json:"title"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

type messageData struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// Telegram ...
type Telegram struct {
	name         string
	token        string
	chatID       string
	lastUpdateID int
	client       *http.Client
}

// NewTelegram ...
func NewTelegram(name, token, chatID string) *Telegram {
	return &Telegram{
		name:   name,
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}
}

// GetName ...
func (t *Telegram) GetName() string {
	return t.name
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

func (t *Telegram) newRequest(method, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("https://api.telegram.org/bot%s%s", t.token, uri), body)
}

// https://core.telegram.org/bots/api#getupdates
func (t *Telegram) getUpdates(offset int) ([]Update, error) {
	var payload = bytes.NewBufferString(
		fmt.Sprintf("{\"offset\": %d}", offset + 1),
	)
	req, err := t.newRequest("GET", "/getUpdates", payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not valid status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Ok      bool     `json:"ok"`
		Updates []Update `json:"result"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data.Updates, nil
}

// GetNewMessages ...
func (t *Telegram) GetNewMessages() ([]*entities.Message, error) {
	updates, err := t.getUpdates(t.lastUpdateID)
	if err != nil {
		return nil, err
	}

	var messages []*entities.Message
	for _, upd := range updates {
		if upd.Message.Text == "" {
			continue
		}

		if strconv.Itoa(upd.Message.Chat.ID) != t.chatID {
			continue
		}

		t.lastUpdateID = upd.UpdateID

		msg := entities.Message{
			Author: entities.Author{
				Username: upd.Message.From.Username,
			},
			Text: upd.Message.Text,
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

// SendMessage ...
func (t *Telegram) SendMessage(m *entities.Message) error {
	msg := messageData{
		ChatID:    t.chatID,
		Text:      m.Text,
		ParseMode: "Markdown",
	}

	d, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := t.newRequest("POST", "/sendMessage", bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("not valid status code: %d msg: %s", resp.StatusCode, b)
	}

	return nil
}
