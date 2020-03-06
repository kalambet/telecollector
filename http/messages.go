package http

import (
	"log"
	"net/http"

	"github.com/kalambet/telecollector/telecollector"
)

func (s *server) handleMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxVal, ok := r.Context().Value(ContextKeyMessage).(*telecollector.MessageContext)
		if !ok {
			s.respond(w, http.StatusInternalServerError, "Message context is invalid")
			return
		}

		text, err := s.msgService.Save(ctxVal)
		if err != nil {
			log.Printf("server: error saving message: %s", err.Error())
			s.respond(w, http.StatusInternalServerError, "Error saving message")
			return
		}

		if ctxVal.Action == telecollector.ActionSave {
			var bcID int64
			if ctxVal.Message.ReplyToMessage != nil {
				bcID, err = s.bot.ForwardMessage(ctxVal.Message.ReplyToMessage.Chat.ID, ctxVal.Message.ReplyToMessage.ID)
				if err != nil {
					log.Printf("server: error forwarding replied message: %s", err.Error())
					s.respond(w, http.StatusInternalServerError, "Error forwarding message")
					return
				}

				bcID, err = s.bot.ReplyBroadcast(text, bcID)
				if err != nil {
					log.Printf("server: error creating reply broadcast: %s", err.Error())
					s.respond(w, http.StatusInternalServerError, "Error creating reply broadcast")
					return
				}

			} else {
				bcID, err = s.bot.ForwardMessage(ctxVal.Message.Chat.ID, ctxVal.Message.ID)
				if err != nil {
					log.Printf("server: error forwarding message: %s", err.Error())
					s.respond(w, http.StatusInternalServerError, "Error forwarding message")
					return
				}
			}

			err = s.msgService.LogBroadcast(ctxVal.Message, bcID)
			if err != nil {
				log.Printf("server: error saving broadcast: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error saving broadcast")
				return
			}
		} else if ctxVal.Action == telecollector.ActionAppend {
			// Append only for connected messages
			// So we need to
			// 1. remove previous broadcast message
			// 2. forward new one
			// 3. create replay to the forwarded one with text from the first
			bcID, err := s.msgService.FindBroadcast(ctxVal.ConnectedMessageID, ctxVal.Message.Chat.ID)
			if err != nil {
				log.Printf("server: error looking for broadcast message: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error looking for broadcast message")
				return
			}

			err = s.bot.DeleteMessage(bcID)
			if err != nil {
				log.Printf("server: error deleting message: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error deleting message")
				return
			}

			bcID, err = s.bot.ForwardMessage(ctxVal.Message.Chat.ID, ctxVal.Message.ID)
			if err != nil {
				log.Printf("server: error forwarding message: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error forwarding message")
				return
			}

			err = s.msgService.LogBroadcast(ctxVal.Message, bcID)
			if err != nil {
				log.Printf("server: error saving broadcast: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error saving broadcast")
				return
			}

			bcID, err = s.bot.ReplyBroadcast(text, bcID)
			if err != nil {
				log.Printf("server: error creating reply broadcast: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error creating reply broadcast")
				return
			}

			err = s.msgService.LogBroadcast(ctxVal.Message, bcID)
			if err != nil {
				log.Printf("server: error saving broadcast: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error saving broadcast")
				return
			}
		} else if ctxVal.Action == telecollector.ActionEdit {
			bcID, err := s.msgService.FindBroadcast(ctxVal.Message.ID, ctxVal.Message.Chat.ID)
			if err != nil {
				log.Printf("server: error looking for broadcast message: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error looking for broadcast message")
				return
			}

			err = s.bot.EditMessage(bcID, text)
			if err != nil {
				log.Printf("server: error editing message: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error looking for broadcast message")
				return
			}
		}
		s.respond(w, http.StatusOK, "OK")
	}
}
