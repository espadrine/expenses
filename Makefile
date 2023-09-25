GOFILES = */*.go src/sql/*.sql

expense: $(GOFILES)
	go build cmd/expense.go

test: $(GOFILES)
	go test ./src

clean:
	rm expense

.PHONY: test clean
