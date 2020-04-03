SHELL := /bin/zsh
export GO111MODULE=on
export GOOGLE_APPLICATION_CREDENTIALS=./service-key.json

run-publisher:
	go build -o publisher cmd/publisher/main.go
	./publisher

run-receiver:
	go build -o receiver cmd/receiver/main.go
	./receiver