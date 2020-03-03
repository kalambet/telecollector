package http

import (
	"log"
	"net/http"
	"sort"

	"github.com/kalambet/telecollector/telecollector"
)

func (s *server) handleMessage(entry *telecollector.Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if entry != nil {
			save, err := s.entryOfInterest(entry)
			if err != nil {
				log.Printf("server: error checking entry as EOI: %s", err.Error())
				return
			}

			if !save {
				s.respond(w, http.StatusOK, "Message was not saved")
				return
			}

			text, err := s.msgService.Save(entry)
			if err != nil {
				log.Printf("server: error saving message: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error saving message")
				return
			}

			switch entry.Message.Action {
			case telecollector.ActionSave:
				msgID, err := s.bot.SendMessage(text)
				if err != nil {
					log.Printf("server: error sending message: %s", err.Error())
					s.respond(w, http.StatusInternalServerError, "Error sending message")
					return
				}

				err = s.msgService.LogBroadcast(entry.Message.ID, msgID)
				if err != nil {
					log.Printf("server: error logging broadcast: %s", err.Error())
					s.respond(w, http.StatusInternalServerError, "Error logging broadcast")
					return
				}

				break
			case telecollector.ActionAppend:
				msgID, err := s.msgService.FindBroadcast(entry.Message.ID)
				if err != nil {
					log.Printf("server: error searching broadcast: %s", err.Error())
					s.respond(w, http.StatusInternalServerError, "Error searching broadcast")
					return
				}

				err = s.bot.EditMessage(msgID, text)
				if err != nil {
					log.Printf("server: error editing message: %s", err.Error())
					s.respond(w, http.StatusInternalServerError, "Error editing message")
					return
				}
				break
			}

		}

		s.respond(w, http.StatusOK, "OK")
	}
}

func (s *server) entryOfInterest(entry *telecollector.Entry) (bool, error) {
	if len(entry.Message.Text) == 0 {
		return false, nil
	}

	connected, err := s.msgService.CheckConnected(entry)
	if err != nil {
		return false, err
	}

	if connected {
		entry.Message.Action = telecollector.ActionAppend
		entry.Message.ID -= 1
		return true, nil
	}

	// Now we only save entries that have trigger tag
	idx := sort.SearchStrings(entry.Message.Tags, telecollector.TriggerTag)
	return idx < len(entry.Message.Tags), nil
}
