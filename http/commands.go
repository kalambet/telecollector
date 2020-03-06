package http

import (
	"log"
	"net/http"

	"github.com/kalambet/telecollector/telecollector"
)

func (s *server) routeCommand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxVal, ok := r.Context().Value(ContextKeyCommand).(*telecollector.CommandContext)
		if !ok {
			s.respond(w, http.StatusInternalServerError, "Command context is invalid")
			return
		}

		if len(ctxVal.Receiver) != 0 && ctxVal.Receiver != s.bot.GetUsername() {
			s.respond(w, http.StatusOK, "OK")
			return
		}

		switch ctxVal.CommandName {
		case telecollector.CommandFollow:
			s.onlyAdminCommand(s.handleFollow())(w, r)
			return
		case telecollector.CommandUnfollow:
			s.onlyAdminCommand(s.handleUnfollow())(w, r)
			return
		case telecollector.CommandWhoami:
			s.handleWhoami()(w, r)
			return
		}

		s.respond(w, http.StatusOK, "OK")
	}
}

func (s *server) handleFollow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxVal, ok := r.Context().Value(ContextKeyCommand).(*telecollector.CommandContext)
		if !ok {
			s.respond(w, http.StatusInternalServerError, "Command context is invalid")
			return
		}

		err := s.credService.FollowChat(ctxVal.Message.Chat, true)
		if err != nil {
			log.Printf("server: follow chat command error: %s", err.Error())
			s.respond(w, http.StatusInternalServerError, "Can not follow this chat")
			return
		}
		s.respond(w, http.StatusOK, "OK")
	}
}

func (s *server) handleUnfollow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxVal, ok := r.Context().Value(ContextKeyCommand).(*telecollector.CommandContext)
		if !ok {
			s.respond(w, http.StatusInternalServerError, "Command context is invalid")
			return
		}

		err := s.credService.FollowChat(ctxVal.Message.Chat, false)
		if err != nil {
			log.Printf("server: unfollow chat command error: %s", err.Error())
			s.respond(w, http.StatusInternalServerError, "Can not unfollow this chat")
			return
		}

		s.respond(w, http.StatusOK, "OK")
	}
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxVal, ok := r.Context().Value(ContextKeyCommand).(*telecollector.CommandContext)
		if !ok {
			s.respond(w, http.StatusInternalServerError, "Command context is invalid")
			return
		}

		_, err := s.bot.ReplyMessage(
			ctxVal.Message.Author().ComposeWhoAmIMessage(), ctxVal.Message.Chat.ID, ctxVal.Message.ID)
		if err != nil {
			log.Printf("server: error sending `whoami` response: %s", err.Error())
			s.respond(w, http.StatusInternalServerError, "Error sending `whoami` message")
			return
		}
		s.respond(w, http.StatusOK, "OK")
	}
}
