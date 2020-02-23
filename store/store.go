package store

import (
	"github.com/kalambet/telecollector/store/postgres"
	"github.com/kalambet/telecollector/telecollector"
)

func NewMessagesService() (telecollector.MessageService, error) {
	return postgres.NewMessagesService()
}

func NewCrenetialService() (telecollector.CredentialService, error) {
	return postgres.NewCredentialService()
}

func Shutdown() error {
	return postgres.Shutdown()
}
