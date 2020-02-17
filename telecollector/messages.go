package telecollector

import (
	"time"

	"github.com/kalambet/telecollector/telegram"
)

type Message struct {
	ID     int64
	Text   string
	Tags   []string
	Date   time.Time
	Author *Author
	Chat   *Chat
}

type Chat struct {
	ID        int64
	Messenger string
	Name      string
}

type Author struct {
	ID       int64
	First    string
	Last     string
	Username string
}

type MessageService interface {
	Save(*Message) error
}

func NewMessage(upd *telegram.Update) *Message {
	var msg *telegram.Message
	if upd.Message != nil {
		msg = upd.Message
	} else if upd.EditedMessage != nil {
		msg = upd.EditedMessage
	} else if upd.ChannelPost != nil {
		msg = upd.ChannelPost
	} else if upd.EditedChannelPost != nil {
		msg = upd.EditedChannelPost
	}

	if msg == nil {
		return nil
	}

	if len(msg.Entities) == 0 {
		return nil
	}

	tags := make([]string, 0)
	if msg.Entities != nil {
		for _, e := range msg.Entities {
			tags = append(tags, msg.Text[e.Offset:e.Offset+e.Length])
		}
	}

	if len(tags) == 0 {
		return nil
	}

	return &Message{
		ID:   msg.ID,
		Text: msg.Text,
		Tags: tags,
		Date: time.Unix(msg.Date, 0),
		Chat: &Chat{
			ID:        msg.Chat.ID,
			Messenger: "Telegram",
			Name:      msg.Chat.Title,
		},
		Author: &Author{
			ID:       msg.From.ID,
			First:    msg.From.FirstName,
			Last:     msg.From.LastName,
			Username: msg.From.UserName,
		},
	}
}
