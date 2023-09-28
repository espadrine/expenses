GOFILES = */*.go src/sql/*.sql

expense: $(GOFILES)
	go build cmd/expense.go

test: $(GOFILES) ./expense
	go test ./src ./cmd

clean:
	rm expense

.PHONY: test clean
