package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kalambet/telecollector/telecollector"
	"github.com/kalambet/telecollector/telegram"
)

func (s *server) buildContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.respond(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.respond(w, http.StatusBadRequest, "Body can't be read")
		}
		log.Printf("The Body: \n%s", body)

		d := json.NewDecoder(bytes.NewReader(body))
		var upd telegram.Update
		err = d.Decode(&upd)
		if err != nil {
			s.respond(w, http.StatusNotAcceptable, "Sent entity is not update")
			return
		}

		var ctx context.Context
		if r.Context() != nil {
			ctx = r.Context()
		} else {
			ctx = context.Background()
		}

		ctx = context.WithValue(ctx, ContextKeyUpdate, &upd)
		r = r.WithContext(ctx)

		if next != nil {
			next(w, r)
		}
	}
}

func (s *server) onlyAdminCommand(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxVal, ok := r.Context().Value(ContextKeyCommand).(*telecollector.CommandContext)
		if !ok {
			s.respond(w, http.StatusNotAcceptable, "Command context is invalid")
			return
		}

		if !s.credService.CheckAdmin(ctxVal.Message.From.ID) {
			s.respond(w, http.StatusNotAcceptable, "Sent entry is not from authorized admin")
			return
		}

		if next != nil {
			next(w, r)
		}
	}
}

func (s *server) onlyWhitelistedChats(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxVal, ok := r.Context().Value(ContextKeyMessage).(*telecollector.MessageContext)
		if !ok {
			s.respond(w, http.StatusNotAcceptable, "Message context is invalid")
			return
		}

		if !s.credService.CheckChat(ctxVal.Message.Chat.ID) {
			s.respond(w, http.StatusNotAcceptable, "Sent entry is not from followed chat")
			return
		}

		if next != nil {
			next(w, r)
		}
	}
}
