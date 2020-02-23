package telecollector

type Allowance struct {
	ChatID   int64
	AuthorID int64
	Follow   bool
}

type CredentialService interface {
	CheckAdmin(int64) bool
	CheckChat(int64) bool
	FollowChat(*Allowance) error
}
