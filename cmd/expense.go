package main

import (
	"database/sql"
	"fmt"
	"log"

	expense "github.com/espadrine/expenses/src"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	params := expense.ParseFlags()
	db := openDB()
	defer db.Close()
	params.Command.Execute(&params)
	// FIXME: the logic should be inside Execute().
	switch params.Command.Names[0] {
	case "users":
		fallthrough
	case "list":
		listUsers(db)
	}
}

func openDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func listUsers(db *sql.DB) {
	sqlStmt := `select * from users`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}
	defer rows.Close()
	fmt.Println("Users:")
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("-", name)
	}
}
