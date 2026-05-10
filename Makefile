APP_NAME=penetapan-service

.PHONY: build run clean test swagger

build:
	go build -o ./bin/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

swagger:
	swag init -g ./cmd/api/main.go

clean:
	rm -f ./bin/$(APP_NAME)
