package main

import (
	"log"
	"os"

	expense "github.com/espadrine/expenses/src"
)

func main() {
	params := expense.ParseFlags(os.Args[1:])
	store, err := expense.NewDB()
	if err != nil {
		log.Fatalf("NewDB: %s\n", err)
	}
	defer store.Close()
	params.Command.Execute(&params, store)
}
