.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o comments-api main.go

.PHONY: run
run: build
	./comments-api "$@"

.PHONY: dev
dev:
	air
