package telecollector

import (
	"fmt"
	"strings"
	"time"

	"github.com/kalambet/telecollector/telegram"
)

const (
	TriggerTag      = "#a51"
	CommandFollow   = "follow"
	CommandUnfollow = "unfollow"
	CommandWhoami   = "whoami"
)

type Bot interface {
	GetUsername() string
	SendMessage(int64, string) error
}

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
	Name     string
	Receiver string
	Params   map[string]string
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
	var msg *telegram.Message
	if upd.Message != nil {
		msg = upd.Message
	}

	if upd.EditedMessage != nil {
		msg = upd.EditedMessage
	}

	if upd.ChannelPost != nil {
		msg = upd.ChannelPost
	}

	if upd.EditedChannelPost != nil {
		msg = upd.EditedChannelPost
	}

	if msg == nil {
		return nil
	}

	entry := &Entry{
		Message: &Message{
			ID: upd.ID,
		},
		Chat: &Chat{
			ID:        msg.Chat.ID,
			Messenger: "Telegram",
			Name:      msg.Chat.Title,
		},
	}

	if msg.Entities == nil || len(msg.Entities) == 0 {
		return nil
	}

	entry.Message.Tags = make([]string, 0)
	for _, e := range msg.Entities {
		if e.Type == telegram.EntityTypeHashtag {
			entry.Message.Tags = append(entry.Message.Tags, msg.Text[e.Offset:e.Offset+e.Length])
		} else if e.Type == telegram.EntityTypeBotCommand {
			// bot command looks like `/command@NameBot`
			// so we split string by @ and then take first segment from second letter to the end
			parts := strings.Split(msg.Text[e.Offset:e.Offset+e.Length], "@")
			entry.Command = &Command{
				Name:     parts[0][1:],
				Receiver: parts[len(parts)-1],
				Params:   nil,
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

	if !entry.isApplicable() {
		return nil
	}

	if entry.Author == nil {
		// Telegram `From` is empty if the message is from Channel
		entry.Author = &Author{
			ID:    0,
			First: msg.Chat.Title,
		}
	}

	return entry
}

func NewBot(token string) (Bot, error) {
	return telegram.NewBot(token)
}

func ComposeWhoAmIMessage(author *Author) string {
	return fmt.Sprintf(
		"*Name*: %s %s\n*Username*: %s\n*ID*:%d",
		author.First,
		author.Last,
		author.Username,
		author.ID)
}
