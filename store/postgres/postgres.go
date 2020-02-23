package postgres

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("postgres: startup problem: %s", err.Error())
	}
}

func Shutdown() error {
	return db.Close()
}

func gracefulCreateTable(table string, query string) error {
	rows, err := db.Query(queryTableExistence, table)
	if err != nil {
		return nil
	}

	if rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if !exists || err == sql.ErrNoRows {
			_, err := db.Query(query)
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
