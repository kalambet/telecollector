package postgres

import "github.com/kalambet/telecollector/telecollector"

type Service struct{}

func NewMessagesService() telecollector.MessageService {
	return &Service{}
}

func (s *Service) Save(msg *telecollector.Message) error {
	return nil
}
