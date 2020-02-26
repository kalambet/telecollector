package http

import (
	"fmt"
	"net/http"

	"github.com/kalambet/telecollector/telecollector"
)

// The only way to 'secure' the endpoint for telegram bot is to make it 'un-discoverable'
//
// Quote from the official Telegram doc: https://core.telegram.org/bots/api#setwebhook
//
// > If you'd like to make sure that the Webhook request comes from Telegram, we recommend
// > using a secret path in the URL, e.g. https://bot.example.com/<token>. Since nobody
// > else knows your bot‘s token, you can be pretty sure it’s us.

func (s *server) routes(secretPath string) {
	s.router.HandleFunc("/", s.handleStatus())
	s.router.HandleFunc(fmt.Sprintf("/%s", secretPath), s.buildContext(s.routeUpdate()))
}

func (s *server) routeUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entry, ok := r.Context().Value(ContextKeyEntry).(*telecollector.Entry)
		if !ok {
			s.respond(w, http.StatusNotAcceptable, "Unable to rote update")
			return
		}

		if entry.Command != nil {
			s.routeCommand(entry)(w, r)
		} else {
			s.onlyWhitelistedChats(s.handleMessage(entry))(w, r)
		}
	}
}

func (s *server) routeCommand(entry *telecollector.Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(entry.Command.Receiver) != 0 && entry.Command.Receiver != s.bot.GetUsername() {
			s.respond(w, http.StatusOK, "OK")
			return
		}

		switch entry.Command.Name {
		case telecollector.CommandFollow:
			s.onlyAdmin(s.handleFollow(entry, true))(w, r)
			return
		case telecollector.CommandUnfollow:
			s.onlyAdmin(s.handleFollow(entry, false))(w, r)
			return
		case telecollector.CommandWhoami:
			s.handleWhoami(entry)(w, r)
			return
		}

		s.respond(w, http.StatusOK, "OK")
	}
}
