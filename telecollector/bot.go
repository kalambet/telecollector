package telecollector

import "github.com/kalambet/telecollector/telegram"

type Bot interface {
	GetUsername() string
	SendMessage(string) (int64, error)
	EditMessage(int64, string) error
}

func NewBot(token string) (Bot, error) {
	return telegram.NewBot(token)
}
