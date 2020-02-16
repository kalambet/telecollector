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

	if upd.Message != nil {
		log.Printf("Message recived from telegram bot: \n%#v", *upd.Message)
	}

	if upd.EditedMessage != nil {
		log.Printf("EditedMessage recived from telegram bot: \n%#v", *upd.EditedMessage)
	}

	if upd.ChannelPost != nil {
		log.Printf("ChannelPost recived from telegram bot: \n%#v", *upd.ChannelPost)
	}

	if upd.EditedChannelPost != nil {
		log.Printf("EditedChannelPost recived from telegram bot: \n%#v", *upd.EditedChannelPost)
	}

	return &Message{}
}
