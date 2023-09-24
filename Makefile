expense: */*.go src/sql/*.sql
	go build cmd/expense.go

clean:
	rm expense

.PHONY: clean
