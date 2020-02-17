package telecollector

import (
	"time"

	"github.com/kalambet/telecollector/telegram"
)

type Message struct {
	ID       int64     `sql:"message_id"`
	Text     string    `sql:"text"`
	Tags     []string  `sql:"tags"`
	Date     time.Time `sql:"date"`
	AuthorID int64     `sql:"author_id"`
	ChatID   int64     `sql:"chat_id"`
}

type Chat struct {
	ID        int64  `sql:"author_id"`
	Messenger string `sql:"messenger"`
	Name      string `sql:"name"`
}

type Author struct {
	ID       int64  `sql:"chat_id"`
	First    string `sql:"first"`
	Last     string `sql:"last"`
	Username string `sql:"username"`
}

type Entry struct {
	Message Message
	Author  Author
	Chat    Chat
}

type MessageService interface {
	Save(*Entry) error
}

func NewEntry(upd *telegram.Update) *Entry {
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

	author := Author{}
	if msg.From != nil {
		author = Author{
			ID:       msg.From.ID,
			First:    msg.From.FirstName,
			Last:     msg.From.LastName,
			Username: msg.From.UserName,
		}
	} else {
		// Telegram `From` is empty if the message is from Channel
		author = Author{
			ID:    0,
			First: msg.Chat.Title,
		}
	}

	return &Entry{
		Message: Message{
			ID:       msg.ID,
			Text:     msg.Text,
			Tags:     tags,
			Date:     time.Unix(msg.Date, 0),
			ChatID:   msg.Chat.ID,
			AuthorID: author.ID,
		},
		Chat: Chat{
			ID:        msg.Chat.ID,
			Messenger: "Telegram",
			Name:      msg.Chat.Title,
		},
		Author: author,
	}
}
