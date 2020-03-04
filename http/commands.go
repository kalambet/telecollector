package http

import (
	"log"
	"net/http"

	"github.com/kalambet/telecollector/telecollector"
)

func (s *server) handleFollow(entry *telecollector.Entry, direction bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.credService.FollowChat(&telecollector.Allowance{
			ChatID:   entry.Chat.ID,
			AuthorID: entry.Author.ID,
			Follow:   direction,
		})
		if err != nil {
			log.Printf("server: follow/unfollow chat command error: %s", err.Error())
			s.respond(w, http.StatusInternalServerError, "Can not follow/unfollow this chat")
			return
		}

	}
}

func (s *server) handleWhoami(entry *telecollector.Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.bot.ReplyMessage(telecollector.ComposeWhoAmIMessage(entry.Author), entry.Chat.ID, entry.Message.Nonce)
		if err != nil {
			log.Printf("server: error sending `whoami` response: %s", err.Error())
			s.respond(w, http.StatusInternalServerError, "Error sending `whoami` message")
			return
		}
	}
}
