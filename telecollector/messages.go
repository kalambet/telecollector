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
	log.Printf("Update recived from telegram bot:\n%#v\n\n", *upd)

	var msg *telegram.Message
	if upd.Message != nil {
		log.Printf("Message recived from telegram bot:\n%#v", *upd.Message)
		msg = upd.Message
	}
	if upd.EditedMessage != nil {
		log.Printf("EditedMessage recived from telegram bot:\n%#v", *upd.EditedMessage)
		msg = upd.EditedMessage
	}
	if upd.ChannelPost != nil {
		log.Printf("ChannelPost recived from telegram bot:\n%#v", *upd.ChannelPost)
		msg = upd.ChannelPost
	}
	if upd.EditedChannelPost != nil {
		log.Printf("EditedChannelPost recived from telegram bot:\n%#v", *upd.EditedChannelPost)
		msg = upd.EditedChannelPost
	}
	if msg != nil {
		log.Printf("Chat recived from telegram bot:\n%#v", *msg.Chat)

		if msg.Entities != nil {
			for _, e := range msg.Entities {
				log.Printf("Entity: %#v\n", e)
			}
		}
	}

	return &Message{}
}
