APP_NAME=salesforge
CMD_DIR=cmd
MAIN_FILE=main.go
PKG=./...

.PHONY: run test fmt lint

run:
	go run $(MAIN_FILE)

test:
	go test -v $(PKG)

fmt:
	go fmt $(PKG)

lint:
	golangci-lint run

