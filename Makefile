GOFILES = */*.go src/sql/*.sql

expense: $(GOFILES)
	go build cmd/expense.go

test: $(GOFILES) ./expense
	go test ./src ./cmd
	rm cmd/.test*.sqlite

clean:
	rm expense

.PHONY: test clean
