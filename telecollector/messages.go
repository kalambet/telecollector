package telecollector

import (
	"time"

	"github.com/kalambet/telecollector/telegram"
)

type Message struct {
	ID       int64
	ChatID   int64
	AuthorID int64
	Date     time.Time
	Text     string
	Tags     []string
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

type Command struct {
	Name   string
	Params map[string]string
}

type Entry struct {
	Message *Message
	Author  *Author
	Chat    *Chat
	Command *Command
}

func (e *Entry) isApplicable() bool {
	return (e.Message != nil && len(e.Message.Tags) > 0) || e.Command != nil
}

type MessageService interface {
	Save(*Entry) error
}

func NewEntry(upd *telegram.Update) *Entry {
	msgs := make([]*telegram.Message, 0)
	if upd.Message != nil {
		msgs = append(msgs, upd.Message)
	}

	if upd.EditedMessage != nil {
		msgs = append(msgs, upd.EditedMessage)
	}

	if upd.ChannelPost != nil {
		msgs = append(msgs, upd.ChannelPost)
	}

	if upd.EditedChannelPost != nil {
		msgs = append(msgs, upd.EditedChannelPost)
	}

	if len(msgs) == 0 {
		return nil
	}

	entry := &Entry{
		Message: &Message{
			ID: upd.ID,
		},
		Chat: &Chat{},
	}

	for _, msg := range msgs {
		if msg.Entities == nil || len(msg.Entities) == 0 {
			continue
		}

		entry.Message.Tags = make([]string, 0)
		for _, e := range msg.Entities {
			if e.Type == "hashtag" {
				entry.Message.Tags = append(entry.Message.Tags, msg.Text[e.Offset:e.Offset+e.Length])
			} else if e.Type == "bot_command" {
				entry.Command = &Command{
					Name:   msg.Text[e.Offset : e.Offset+e.Length],
					Params: nil,
				}
			}
		}

		if msg.From != nil {
			entry.Author = &Author{
				ID:       msg.From.ID,
				First:    msg.From.FirstName,
				Last:     msg.From.LastName,
				Username: msg.From.UserName,
			}
		}
	}

	if !entry.isApplicable() {
		return nil
	}

	if entry.Author == nil {
		// Telegram `From` is empty if the message is from Channel
		entry.Author = &Author{
			ID:    0,
			First: msgs[0].Chat.Title,
		}
	}

	return entry
}
