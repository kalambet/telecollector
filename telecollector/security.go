package telecollector

import "github.com/kalambet/telecollector/telegram"

type Allowance struct {
	ChatID int64
	Follow bool
}

type CredentialService interface {
	CheckAdmin(int64) bool
	CheckChat(int64) bool
	FollowChat(*telegram.Chat, bool) error
}
