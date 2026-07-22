APP_NAME=penetapan-service

.PHONY: build run clean test swagger build-image build-docker format dev

build:
	@mkdir -p bin
	go build -o ./bin/$(APP_NAME) ./cmd/api

build-docker:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" \
	-o ./bin/$(APP_NAME) ./cmd/api

format:
	goimports -w .

run:
	go run ./cmd/api

dev: format run

test:
	go test -race ./...

swagger:
	swag init -d cmd/api,internal/api,internal/model/domain,internal/model/web,internal/individu/domain,internal/individu/web,internal/pemda/domain,internal/pemda/web

clean:
	rm -f ./bin/$(APP_NAME)

build-image:
	@docker build . -t $(APP_NAME)
