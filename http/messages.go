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
			// We only save entries that have trigger tag
			idx := sort.SearchStrings(entry.Message.Tags, telecollector.TriggerTag)
			if idx < len(entry.Message.Tags) {
				err := s.msgService.Save(entry)
				if err != nil {
					log.Printf("server: error saving message: %s", err.Error())
					s.respond(w, http.StatusInternalServerError, "Error saving message")
					return
				}

				if err = s.bot.RepostMessage(entry.Message.Text); err != nil {
					log.Printf("server: error reposting message: %s", err.Error())
				}
			}
		}

		s.respond(w, http.StatusOK, "OK")
	}
}
