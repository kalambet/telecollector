package postgres

import (
	"context"
	"database/sql"
	"os"

	"github.com/kalambet/telecollector/telecollector"
	_ "github.com/lib/pq"
)

const (
	messageTableQuery = `
create table messages(
    message_id integer, 
    chat_id integer, 
    date date, 
    text text not null, 
    tags text[] not null default '{}', 
    author_id int, 
    primary key(chat_id, message_id)
);`

	authorsTableQuery = `
create table authors(
    author_id integer primary key, 
    first text not null, 
    last text, 
    username text
);`

	chatsTableQuery = `
create table chats(
    chat_id integer primary key, 
    messenger text not null, 
    name text 
);`

	tableExistenceSQL = `select exists (select from pg_tables where schemaname = 'public' and tablename = $1);`

	insertChatQuery    = `insert into chats (chat_id, messenger, name) values (:chat_id, :messenger, :name) on conflict (chat_id) do nothing;`
	insertAuthorQuery  = `insert into authors (author_id, first, last, username) values(:author_id, :first, :last, :username) on conflict (author_id) do nothing;`
	insertMessageQuery = `
insert into 
    messages (message_id, chat_id, date, text, tags, author_id) 
    values (:message_id, :chat_id, :date, :text, :tags, :author_id) 
    on conflict do update set date = :date, text = :text, tags = :tags;`
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
	tx, err := s.db.BeginTx(context.Background(), &sql.TxOptions{
		ReadOnly: false,
	})

	_, err = tx.Exec(insertChatQuery, entity.Chat)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	_, err = tx.Exec(insertAuthorQuery, entity.Author)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	_, err = tx.Exec(insertMessageQuery, entity.Message)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
