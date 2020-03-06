package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/kalambet/telecollector/telegram"

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
		upd, ok := r.Context().Value(ContextKeyUpdate).(*telegram.Update)
		if !ok {
			s.respond(w, http.StatusNotAcceptable, "Unable to rote update")
			return
		}

		action := telecollector.ActionSave
		var msg *telegram.Message
		if upd.Message != nil {
			msg = upd.Message
		}

		if upd.EditedMessage != nil {
			msg = upd.EditedMessage
			action = telecollector.ActionEdit
		}

		if upd.ChannelPost != nil {
			msg = upd.ChannelPost
		}

		if upd.EditedChannelPost != nil {
			msg = upd.EditedChannelPost
			action = telecollector.ActionEdit
		}

		if msg == nil {
			s.respond(w, http.StatusNotAcceptable, "Update has no message")
			return
		}

		cmd, rcvr := msg.Command()
		if len(cmd) != 0 {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKeyCommand, &telecollector.CommandContext{
				Message:      msg,
				CommandName:  cmd,
				CommandPrams: nil,
				Receiver:     rcvr,
			})

			s.routeCommand()(w, r.WithContext(ctx))
			return
		}

		connected, err := s.msgService.CheckConnected(msg)
		if err != nil {
			log.Printf("server: error checking entry as UOI: %s", err.Error())
			return
		}

		if connected {
			action = telecollector.ActionAppend
		} else {
			tags := msg.Tags()
			if sort.SearchStrings(tags, telecollector.TriggerTag) == len(tags) {
				s.respond(w, http.StatusOK, "OK")
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyUpdate, nil)
		ctx = context.WithValue(ctx, ContextKeyMessage, &telecollector.MessageContext{
			Message:            msg,
			ConnectedMessageID: msg.ID - 1,
			UpdateID:           upd.ID,
			Action:             action,
		})
		s.onlyWhitelistedChats(s.handleMessage())(w, r.WithContext(ctx))
	}
}
