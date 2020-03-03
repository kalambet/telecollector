package telecollector

import "github.com/kalambet/telecollector/telegram"

type Bot interface {
	GetUsername() string
	SendMessage(string) (int64, error)
	EditMessage(int64, string) error
	ForwardMessage(int64, int64) error
}

func NewBot(token string) (Bot, error) {
	return telegram.NewBot(token)
}
