package telecollector

type Message struct{}

type Update struct {
	UpdaterId int       `json:"update_id"`
	Message   TGMessage `json:"message"`
}

type MessageService interface {
	Save(*Message) error
}

func (u *Update) Message() *Message {
	return nil
}
