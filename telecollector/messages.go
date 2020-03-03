package telecollector

import (
	"fmt"
	"strings"

	"github.com/kalambet/telecollector/telegram"
)

const (
	TriggerTag      = "#a51"
	CommandFollow   = "follow"
	CommandUnfollow = "unfollow"
	CommandWhoami   = "whoami"
)

type MessageAction string

var (
	ActionAppend MessageAction = "append"
	ActionSave   MessageAction = "save"
)

type Message struct {
	ID       int64
	Nonce    int64
	ChatID   int64
	AuthorID int64
	Date     int64
	Text     string
	Tags     []string
	Action   MessageAction
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
	Chat     *Chat
}

type Entry struct {
	Message *Message
	Author  *Author
	Chat    *Chat
	Command *Command
}

type MessageService interface {
	Save(*Entry) (string, error)
	LogBroadcast(int64, int64) error
	FindBroadcast(int64) (int64, error)
	CheckConnected(*Entry) (bool, error)
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
			ID:     upd.ID,
			Text:   msg.Text,
			Date:   msg.Date,
			Action: ActionSave,
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
			parts := strings.Split(upd.Message.Text[e.Offset:e.Offset+e.Length], "@")
			var receiver string
			if len(parts) == 1 {
				receiver = ""
			} else {
				receiver = parts[len(parts)-1]
			}
			entry.Command = &Command{
				Name:     parts[0][1:],
				Receiver: receiver,
				Params: map[string]string{
					"0": upd.Message.Text[e.Offset+e.Length:],
				},
			}
			break
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

	if entry.Author == nil {
		// Telegram `From` is empty if the message is from Channel
		entry.Author = &Author{
			ID:    0,
			First: msg.Chat.Title,
		}
	}

	return entry
}

func ComposeWhoAmIMessage(author *Author) string {
	return fmt.Sprintf(
		"*Name*: %s %s\n*Username*: %s\n*ID*:%d",
		author.First,
		author.Last,
		author.Username,
		author.ID)
}
