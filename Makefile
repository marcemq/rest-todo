# Simple Makefile for rest-todo project

EXECNAME = rest-todo-exec

build:
	go build -o $(EXECNAME)

test:
	go test -v ./...

coverage:
	go test -v -cover ./...

clean:
	go clean && \
	rm -f $(EXECNAME)
