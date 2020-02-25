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
// > using a secret path in the URL, e.g. https://www.example.com/<token>. Since nobody
// > else knows your bot‘s token, you can be pretty sure it’s us.

func (s *server) routes(secretPath string) {
	s.router.HandleFunc("/", s.handleStatus())
	s.router.HandleFunc(fmt.Sprintf("/%s", secretPath), s.buildContext(s.handleMessage()))
}

func (s *server) routeUpdate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(ContextKeyEntry).(*telecollector.Entry)
		if !ok {
			s.respond(w, http.StatusNotAcceptable, "Unable to rote update")
			return
		}

	}
}
