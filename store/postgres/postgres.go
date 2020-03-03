package postgres

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

const queryTableExistence = `select exists (select from pg_tables where schemaname = 'public' and tablename = $1);`

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
	defer rows.Close()

	if err != nil {
		return err
	}

	if rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if !exists || err == sql.ErrNoRows {
			rows, err := db.Query(query)
			defer rows.Close()

			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}
