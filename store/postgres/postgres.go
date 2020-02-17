package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"

	"github.com/kalambet/telecollector/telecollector"
	_ "github.com/lib/pq"
)

const (
	messageTableQuery = `
create table messages(
    message_id bigint, 
    chat_id bigint, 
    author_id bigint,
    date date, 
    text text not null, 
    tags text[] not null default '{}', 
    primary key(chat_id, message_id)
);`

	authorsTableQuery = `
create table authors(
    author_id bigint primary key, 
    first text not null, 
    last text, 
    username text
);`

	chatsTableQuery = `
create table chats(
    chat_id bigint primary key, 
    messenger text not null, 
    name text 
);`

	tableExistenceSQL = `select exists (select from pg_tables where schemaname = 'public' and tablename = $1);`

	insertChatQuery = `
insert into 
    chats (chat_id, messenger, name) 
    values ($1, $2, $3) 
    on conflict (chat_id) do nothing;`

	insertAuthorQuery = `
insert into 
    authors (author_id, first, last, username) 
    values($1, $2, $3, $4) 
    on conflict (author_id) do nothing;`

	insertMessageQuery = `
insert into 
    messages (message_id, chat_id, author_id, date, text, tags) 
    values ($1, $2, $3, $4, $5, $6) 
    on conflict (message_id, chat_id) 
        do update set date = $4, text = $5, tags = $6;`
)

type Service struct {
	db *sql.DB
}

func NewMessagesService() (telecollector.MessageService, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	s := &Service{
		db: db,
	}

	err = s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) init() error {
	err := s.gracefulCreateTable("messages", messageTableQuery)
	if err != nil {
		return err
	}

	err = s.gracefulCreateTable("authors", authorsTableQuery)
	if err != nil {
		return err
	}

	err = s.gracefulCreateTable("chats", chatsTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) gracefulCreateTable(table string, query string) error {
	rows, err := s.db.Query(tableExistenceSQL, table)
	if err != nil {
		return nil
	}

	if rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if !exists || err == sql.ErrNoRows {
			_, err := s.db.Query(query)
			if err != nil {
				return err
			}
		}
		err = rows.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Save(entity *telecollector.Entry) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(insertChatQuery, entity.Chat.ID, entity.Chat.Messenger, entity.Chat.Name)

	if err != nil {
		return rollback(tx, err)
	}

	_, err = tx.Exec(insertAuthorQuery,
		entity.Author.ID, entity.Author.First, entity.Author.Last, entity.Author.Username)

	if err != nil {
		return rollback(tx, err)
	}

	_, err = tx.Exec(insertMessageQuery,
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
