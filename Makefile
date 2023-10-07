GOFILES = */*.go src/sql/*.sql

expense: $(GOFILES)
	go build cmd/expense.go

test: $(GOFILES) ./expense
	rm -f cmd/.test*.sqlite
	go test ./src ./cmd
	rm -f cmd/.test*.sqlite

clean:
	rm expense

.PHONY: test clean
