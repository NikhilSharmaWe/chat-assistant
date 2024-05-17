build:
	go build -o bin/chat-assistant

run: build
	./bin/chat-assistant
