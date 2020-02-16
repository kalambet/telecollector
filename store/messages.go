package store

import (
	"github.com/kalambet/telecollector/store/postgres"
	"github.com/kalambet/telecollector/telecollector"
)

func NewMessagesService() (telecollector.MessageService, error) {
	return postgres.NewMessagesService(), nil
}
