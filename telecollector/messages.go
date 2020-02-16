package telecollector

import (
	"log"

	"github.com/kalambet/telecollector/telegram"
)

type Message struct{}

type MessageService interface {
	Save(*Message) error
}

func NewMessage(upd *telegram.Update) *Message {
	log.Printf("Update recived from telegram bot: \n%#v\n\n", *upd)
	log.Printf("Message recived from telegram bot: \n%#v", *upd.Message)
	return &Message{}
}
