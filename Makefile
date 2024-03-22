version = development-0.0.0

build:
	go build -o comments-api main.go

run: build
	./comments-api "$@"

dev:
	air
