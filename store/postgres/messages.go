package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kalambet/telecollector/telegram"

	"github.com/kalambet/telecollector/telecollector"
	"github.com/lib/pq"
)

const (
	createMessage = `
create table messages(
    update_id bigint unique,
	message_id bigint,
    chat_id bigint, 
    author_id bigint,
    date bigint, 
    text text not null, 
    tags text[] not null default '{}', 
    primary key(message_id, chat_id)
);`

	createAuthors = `
create table authors(
    author_id bigint primary key, 
    first text not null, 
    last text, 
    username text
);`

	createChats = `
create table chats(
    chat_id bigint primary key, 
    messenger text not null, 
    name text 
);`

	createBroadcast = `
create table broadcasts(
    message_id bigint,
	chat_id bigint,
    broadcast_id bigint,
	primary key(message_id, chat_id)
);`

	queryMessagesExistence = `
select exists (select from messages where message_id = $1 and chat_id = $2 and author_id = $3 and date = $4);`

	insertChat = `
insert into 
    chats (chat_id, messenger, name) 
    values ($1, $2, $3) 
    on conflict (chat_id) do nothing;`

	insertAuthor = `
insert into 
    authors (author_id, first, last, username) 
    values($1, $2, $3, $4) 
    on conflict (author_id) do nothing;`

	insertMessage = `
insert into 
    messages (update_id, message_id, chat_id, author_id, date, text, tags) 
    values ($1, $2, $3, $4, $5, $6, $7) 
    on conflict (message_id, chat_id) 
        do update set date = $5, text = $6, tags = $7 returning text;`

	appendMessage = `
insert into 
    messages (update_id, message_id, chat_id, author_id, date, text, tags) 
    values ($1, $2, $3, $4, $5, $6, $7) 
    on conflict (message_id, chat_id) 
        do update set date = $5, text = messages.text || ' ➜ ' || $6 returning text;`

	insertBroadcast = `
insert into
	broadcasts (message_id, chat_id, broadcast_id)
	values ($1, $2, $3)
	on conflict (message_id, chat_id)
		do update set broadcast_id = $3;`

	queryBroadcast = `select broadcast_id from broadcasts where message_id = $1 and chat_id = $2;`
)

type messagesService struct{}

func NewMessagesService() (*messagesService, error) {

	err := gracefulCreateTable("messages", createMessage)
	if err != nil {
		return nil, err
	}

	err = gracefulCreateTable("authors", createAuthors)
	if err != nil {
		return nil, err
	}

	err = gracefulCreateTable("chats", createChats)
	if err != nil {
		return nil, err
	}

	err = gracefulCreateTable("broadcasts", createBroadcast)
	if err != nil {
		return nil, err
	}

	return &messagesService{}, nil
}

func (s *messagesService) Save(ctx *telecollector.MessageContext) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(insertChat, ctx.Message.Chat.ID, "Telegram", ctx.Message.Chat.Title)

	if err != nil {
		return "", rollback(tx, err)
	}

	_, err = tx.Exec(insertAuthor, ctx.Message.From.ID, ctx.Message.From.FirstName, ctx.Message.From.LastName, ctx.Message.From.UserName)

	if err != nil {
		return "", rollback(tx, err)
	}

	var query string
	var msgID int64
	if ctx.Action == telecollector.ActionAppend {
		query = appendMessage
		msgID = ctx.ConnectedMessageID
	} else {
		query = insertMessage
		msgID = ctx.Message.ID
	}

	rows, err := tx.Query(query,
		ctx.UpdateID, msgID, ctx.Message.Chat.ID, ctx.Message.Author().ID,
		ctx.Message.Date, ctx.Message.Text2Save(), pq.Array(ctx.Message.Tags()))

	if err != nil {
		return "", rollback(tx, err)
	}

	var text string
	if rows.Next() {
		err = rows.Scan(&text)
		if err != nil {
			return "", err
		}
	}

	err = rows.Close()
	if err != nil {
		log.Printf("postgres: error closing query rows: %s", err.Error())
	}

	return text, tx.Commit()
}

func (s *messagesService) CheckConnected(msg *telegram.Message) (bool, error) {
	rows, err := db.Query(queryMessagesExistence, msg.ID-1, msg.Chat.ID, msg.Author().ID, msg.Date)

	if err != nil {
		return false, err
	}

	if rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if err != nil {
			return false, nil
		}

		return exists, nil
	}

	err = rows.Close()
	if err != nil {
		return false, err
	}

	return false, nil
}

func (s *messagesService) LogBroadcast(msg *telegram.Message, bcID int64) error {
	_, err := db.Exec(insertBroadcast, msg.ID, msg.Chat.ID, bcID)
	if err != nil {
		return err
	}

	return nil
}

func (s *messagesService) FindBroadcast(msgID int64, chatID int64) (int64, error) {
	rows, err := db.Query(queryBroadcast, msgID, chatID)

	if err != nil {
		return 0, err
	}

	if rows.Next() {
		var bcID int64
		err = rows.Scan(&bcID)
		if err != nil {
			return 0, nil
		}

		return bcID, nil
	}

	err = rows.Close()
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func rollback(tx *sql.Tx, err error) error {
	txErr := tx.Rollback()
	if txErr != nil {
		return fmt.Errorf("%s: %s", err.Error(), txErr.Error())
	}
	return err
}
