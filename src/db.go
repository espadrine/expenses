package expense

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewDB() Store {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	return Store{
		db: db,
	}
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) getUsers() []string {
	var users []string
	sqlStmt := `select * from users`
	rows, err := s.db.Query(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		users = append(users, name)
		if err != nil {
			log.Fatal(err)
		}
	}
	return users
}
