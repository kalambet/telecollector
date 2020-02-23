package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kalambet/telecollector/telecollector"
	"github.com/kalambet/telecollector/telegram"
)

const (
	ContextKeyEntry = "entry_context"
)

func (s *server) buildContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.respond(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		d := json.NewDecoder(r.Body)
		var upd telegram.Update
		err := d.Decode(&upd)
		if err != nil {
			s.respond(w, http.StatusNotAcceptable, "Sent entity is not update")
			return
		}

		entry := telecollector.NewEntry(&upd)
		if entry != nil {
			var ctx context.Context
			if r.Context() != nil {
				ctx = r.Context()
			} else {
				ctx = context.Background()
			}

			ctx = context.WithValue(ctx, ContextKeyEntry, &entry)
			r = r.WithContext(ctx)
		}

		if next != nil {
			next(w, r)
		}
	}
}

func (s *server) onlyWhitelisted(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entry, ok := r.Context().Value(ContextKeyEntry).(*telecollector.Entry)
		if !ok {
			s.respond(w, http.StatusNotAcceptable, "Sent entity is not entry")
			return
		}

		if !s.credService.CheckAdmin(entry.Author.ID) && !s.credService.CheckChat(entry.Chat.ID) {
			s.respond(w, http.StatusNotAcceptable, "Sent entry is not from admin nor from followed chat")
			return
		}

		if next != nil {
			next(w, r)
		}
	}
}
