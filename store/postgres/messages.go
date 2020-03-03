package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kalambet/telecollector/telecollector"
	"github.com/lib/pq"
)

const (
	createMessage = `
create table messages(
    message_id bigint,
	nonce int,
    chat_id bigint, 
    author_id bigint,
    date bigint, 
    text text not null, 
    tags text[] not null default '{}', 
    primary key(chat_id, message_id)
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
    update_id bigint primary key, 
    message_id bigint 
);`

	queryMessagesExistence = `select exists (select from messages where message_id = $1 and chat_id = $2 and author_id = $3 and date = $4);`

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
    messages (message_id, nonce, chat_id, author_id, date, text, tags) 
    values ($1, $2, $3, $4, $5, $6, $7) 
    on conflict (message_id, chat_id) 
        do nothing returning text;`

	appendMessage = `
insert into 
    messages (message_id, nonce, chat_id, author_id, date, text, tags) 
    values ($1, $2, $3, $4, $5, $6, $7) 
    on conflict (message_id, chat_id) 
        do update set date = $5, text = messages.text || ' ➜ ' || $6, tags = $7 returning text;`

	insertBroadcast = `
insert into
	broadcasts (update_id, message_id)
	values ($1, $2)
	on conflict (update_id)
		do update set message_id = $2;
`

	queryBroadcast = `select message_id from broadcasts where update_id = $1;`
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

func (s *messagesService) Save(entity *telecollector.Entry) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(insertChat, entity.Chat.ID, entity.Chat.Messenger, entity.Chat.Name)

	if err != nil {
		return "", rollback(tx, err)
	}

	_, err = tx.Exec(insertAuthor,
		entity.Author.ID, entity.Author.First, entity.Author.Last, entity.Author.Username)

	if err != nil {
		return "", rollback(tx, err)
	}

	var rows *sql.Rows
	if entity.Message.Action == telecollector.ActionAppend {
		rows, err = tx.Query(appendMessage,
			entity.Message.ID, entity.Message.Nonce, entity.Chat.ID, entity.Author.ID,
			entity.Message.Date, entity.Message.Text, pq.Array(entity.Message.Tags))
	} else {
		rows, err = tx.Query(insertMessage,
			entity.Message.ID, entity.Message.Nonce, entity.Chat.ID, entity.Author.ID,
			entity.Message.Date, entity.Message.Text, pq.Array(entity.Message.Tags))
	}

	if err != nil {
		return "", rollback(tx, err)
	}

	var text string
	if rows.Next() {
		err = rows.Scan(&text)
		if err != nil {
			log.Printf("postgres: erorr decoding resulting text: %s", err.Error())
			text = ""
		}
	}

	err = rows.Close()
	if err != nil {
		log.Printf("postgres: error closing query rows: %s", err.Error())
	}

	return text, tx.Commit()
}

func (s *messagesService) CheckConnected(entry *telecollector.Entry) (bool, error) {
	rows, err := db.Query(queryMessagesExistence, entry.Message.ID-1, entry.Chat.ID, entry.Author.ID, entry.Message.Date)
	defer rows.Close()

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

	return false, nil
}

func (s *messagesService) LogBroadcast(updateID int64, messageID int64) error {
	_, err := db.Exec(insertBroadcast, updateID, messageID)
	if err != nil {
		return err
	}

	return nil
}

func (s *messagesService) FindBroadcast(updateID int64) (int64, error) {
	rows, err := db.Query(queryBroadcast, updateID)
	defer rows.Close()

	if err != nil {
		return 0, err
	}

	if rows.Next() {
		var msgID int64
		err = rows.Scan(&msgID)
		if err != nil {
			return 0, nil
		}

		return msgID, nil
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
