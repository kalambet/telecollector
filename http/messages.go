package http

import (
	"log"
	"net/http"

	"github.com/kalambet/telecollector/telecollector"
)

func (s *server) handleMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entry, ok := r.Context().Value(ContextKeyEntry).(*telecollector.Entry)
		if !ok {
			s.respond(w, http.StatusNotAcceptable, "Saving message is not entry")
			return
		}

		if entry != nil {
			err := s.msgService.Save(entry)
			if err != nil {
				log.Printf("server: error saving message: %s", err.Error())
				s.respond(w, http.StatusInternalServerError, "Error saving message")
				return
			}
		}

		s.respond(w, http.StatusOK, "OK")
	}
}
