package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	CommandToMethod = map[string]string{
		"getMe":           http.MethodGet,
		"sendMessage":     http.MethodPost,
		"editMessageText": http.MethodPost,
	}
)

type Bot struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"first_name"`
	channel  int64
	token    string
}

func apiRequest(token string, cmd string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, cmd)
	req, err := http.NewRequest(CommandToMethod[cmd], url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(resp.Body)
	var data Response
	err = d.Decode(&data)
	if err != nil {
		return nil, err
	}

	if !data.OK {
		log.Printf("Result: %s", data.Result)
	}

	return data.Result, nil
}

func NewBot(token string) (*Bot, error) {
	body, err := apiRequest(token, "getMe", nil)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(bytes.NewReader(body))
	var bot Bot
	err = d.Decode(&bot)
	if err != nil {
		return nil, err
	}

	bot.channel, err = strconv.ParseInt(os.Getenv("TG_CHANNEL"), 10, 64)
	if err != nil {
		log.Printf("telegram: error identifing repost channel: %s", err.Error())
		bot.channel = 0
	}

	bot.token = token
	return &bot, nil
}

func (b *Bot) apiRequest(cmd string, body []byte) ([]byte, error) {
	return apiRequest(b.token, cmd, body)
}

func (b *Bot) GetUsername() string {
	return b.Username
}

func (b *Bot) SendMessage(text string) (int64, error) {
	if b.channel == 0 {
		return 0, nil
	}

	msg := struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatId: b.channel,
		Text:   text,
	}

	body, err := json.Marshal(&msg)
	if err != nil {
		return 0, err
	}
	log.Printf("Send Message: %s", body)

	resp, err := b.apiRequest("sendMessage", body)
	if err != nil {
		return 0, err
	}

	respMsg := Message{}
	err = json.Unmarshal(resp, &respMsg)
	if err != nil {
		return 0, nil
	}

	return respMsg.ID, nil
}

func (b *Bot) EditMessage(msgID int64, text string) error {
	if b.channel == 0 {
		return nil
	}

	msg := struct {
		ChatId int64  `json:"chat_id"`
		MsgID  int64  `json:"message_id"`
		Text   string `json:"text"`
	}{
		ChatId: b.channel,
		MsgID:  msgID,
		Text:   text,
	}

	body, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	log.Printf("Send Message: %s", body)

	resp, err := b.apiRequest("editMessageText", body)
	if err != nil {
		return err
	}

	respMsg := Message{}
	err = json.Unmarshal(resp, &respMsg)
	if err != nil {
		return nil
	}

	return nil
}
