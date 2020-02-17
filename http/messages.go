package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kalambet/telecollector/telecollector"

	"github.com/kalambet/telecollector/telegram"
)

func (s *server) handleMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.respond(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		log.Printf("%s", r.Body)

		d := json.NewDecoder(r.Body)
		var upd telegram.Update
		err := d.Decode(&upd)
		if err != nil {
			s.respond(w, http.StatusNotAcceptable, "Sent entity is not update")
			return
		}

		entry := telecollector.NewEntry(&upd)
		if entry != nil {
			err = s.msgService.Save(entry)
			if err != nil {
				log.Printf("server: error saving message: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error saving message")
				return
			}
		}

		s.respond(w, http.StatusOK, "OK")
	}
}
