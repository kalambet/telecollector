package postgres

import (
	"database/sql"
	"fmt"

	"github.com/kalambet/telecollector/telecollector"
	"github.com/lib/pq"
)

const (
	createMessage = `
create table messages(
    message_id bigint, 
    chat_id bigint, 
    author_id bigint,
    date date, 
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

	queryTableExistence = `select exists (select from pg_tables where schemaname = 'public' and tablename = $1);`

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
    messages (message_id, chat_id, author_id, date, text, tags) 
    values ($1, $2, $3, $4, $5, $6) 
    on conflict (message_id, chat_id) 
        do update set date = $4, text = $5, tags = $6;`
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

	return &messagesService{}, nil
}

func (s *messagesService) Save(entity *telecollector.Entry) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(insertChat, entity.Chat.ID, entity.Chat.Messenger, entity.Chat.Name)

	if err != nil {
		return rollback(tx, err)
	}

	_, err = tx.Exec(insertAuthor,
		entity.Author.ID, entity.Author.First, entity.Author.Last, entity.Author.Username)

	if err != nil {
		return rollback(tx, err)
	}

	_, err = tx.Exec(insertMessage,
		entity.Message.ID, entity.Chat.ID, entity.Author.ID,
		entity.Message.Date, entity.Message.Text, pq.Array(entity.Message.Tags))

	if err != nil {
		return rollback(tx, err)
	}

	return tx.Commit()
}

func rollback(tx *sql.Tx, err error) error {
	txErr := tx.Rollback()
	if txErr != nil {
		return fmt.Errorf("%s: %s", err.Error(), txErr.Error())
	}
	return err
}
