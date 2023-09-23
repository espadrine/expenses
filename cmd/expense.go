package main

import (
	expense "github.com/espadrine/expenses/src"
)

func main() {
	params := expense.ParseFlags()
	store := expense.NewDB()
	defer store.Close()
	params.Command.Execute(&params, &store)
}
