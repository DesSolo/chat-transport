package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (t *Telegram) newRequest(method, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("https://api.telegram.org/bot%s%s", t.token, uri), body)
}

// https://core.telegram.org/bots/api#getupdates
func (t *Telegram) getUpdates(offset int) ([]Update, error) {
	var payload = bytes.NewBufferString(
		fmt.Sprintf("{\"offset\": %d}", offset+1),
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



// SendMessage ...
func (t *Telegram) sendMessage(chatID, text string) error {
	msg := messageData{
		ChatID:    chatID,
		Text:      text,
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
