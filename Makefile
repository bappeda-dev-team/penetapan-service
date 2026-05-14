APP_NAME=penetapan-service

.PHONY: build run clean test swagger build-image build-docker

build:
	go build -o ./bin/$(APP_NAME) ./cmd/api

build-docker:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" \
	-o ./bin/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

swagger:
	swag init -d cmd/api,internal/api,internal/model/domain,internal/model/web

clean:
	rm -f ./bin/$(APP_NAME)

build-image:
	@docker build . -t $(APP_NAME)
