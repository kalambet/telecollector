package telecollector

import "github.com/kalambet/telecollector/telegram"

type Bot interface {
	GetUsername() string
	SendMessage(text string) (int64, error)
	EditMessage(msgID int64, text string) error
	ForwardMessage(chatID int64, msgID int64) (int64, error)
	ReplyBroadcast(text string, msgID int64) (int64, error)
	ReplyMessage(text string, chatID int64, msgID int64) (int64, error)
	DeleteMessage(msgID int64) error
}

func NewBot(token string) (Bot, error) {
	return telegram.NewBot(token)
}
