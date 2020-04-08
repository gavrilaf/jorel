SHELL := /bin/zsh
export GO111MODULE=on
export GOOGLE_APPLICATION_CREDENTIALS=./service-key.json

run-publisher:
	go build -o publisher cmd/publisher/main.go
	./publisher

run-publisher2:
	go build -o publisher cmd/publisher/main.go
	./publisher route

run-receiver:
	go build -o receiver cmd/receiver/main.go
	./receiver

run-receiver2:
	go build -o receiver cmd/receiver/main.go
	./receiver cancel-topic-subs


run-jorel:
	go build -o jorel cmd/jorel/main.go
	./jorel