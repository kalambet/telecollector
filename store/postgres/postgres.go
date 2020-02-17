package postgres

import (
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

	insertQuery = `
begin transaction;
insert into chats (chat_id, messenger, name) values ($1, $2, $3) on conflict (chat_id) do nothing;
insert into authors (author_id, first, last, username) values($4, $5, $6, $7) on conflict (author_id) do nothing;
insert into 
    messages (message_id, chat_id, date, text, tags, author_id) 
    values ($8, $1, $9, $10, $11, $4) 
    on conflict do update set date = $9, text = $10, tags = $11;
commit;`
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

func (s *Service) Save(msg *telecollector.Message) error {
	_, err := s.db.Query(insertQuery, msg.Chat.ID, msg.Chat.Messenger, msg.Chat.Name,
		msg.Author.ID, msg.Author.First, msg.Author.Last, msg.Author.Username,
		msg.ID, msg.Date, msg.Text, msg.Tags)
	return err
}
