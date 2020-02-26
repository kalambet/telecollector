package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	CommandToMethod = map[string]string{
		"getMe":       http.MethodGet,
		"sendMessage": http.MethodPost,
	}
)

type Bot struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"first_name"`
	token    string
}

func apiRequest(token string, cmd string, body io.Reader) (io.ReadCloser, error) {
	url := fmt.Sprintf("https://https://api.telegram.org/bot%s/%s", token, cmd)
	req, err := http.NewRequest(CommandToMethod[cmd], url, body)
	if err != nil {
		return nil, err
	}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil

}

func NewBot(token string) (*Bot, error) {
	body, err := apiRequest(token, "getMe", nil)
	if err != nil {
		return nil, err
	}

	var bot Bot
	d := json.NewDecoder(body)
	err = d.Decode(&bot)
	if err != nil {
		return nil, err
	}

	bot.token = token

	return &bot, nil
}

func (b *Bot) apiRequest(cmd string, body io.Reader) (io.ReadCloser, error) {
	return apiRequest(b.token, cmd, body)
}

func (b *Bot) GetUsername() string {
	return b.Username
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := MessageResponse{
		ChatId:    chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}

	body, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	_, err = b.apiRequest("sendMessage", bytes.NewReader(body))
	if err != nil {
		return err
	}

	return nil
}
