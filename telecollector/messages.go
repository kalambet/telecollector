package telecollector

import (
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
	ActionEdit   MessageAction = "edit"
)

type MessageContext struct {
	Message            *telegram.Message
	ConnectedMessageID int64
	UpdateID           int64
	Action             MessageAction
}

type CommandContext struct {
	Message      *telegram.Message
	CommandName  string
	CommandPrams interface{}
	Receiver     string
}

type MessageService interface {
	Save(ctx *MessageContext) (string, error)
	LogBroadcast(msg *telegram.Message, bcID int64) error
	FindBroadcast(msgID int64, chatID int64) (int64, error)
	CheckConnected(msg *telegram.Message) (bool, error)
}
