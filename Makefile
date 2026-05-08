APP_NAME=penetapan-service

.PHONY: all build run myenv clean

# DEFAULT TARGET
all: build

build: $(APP_NAME)

$(APP_NAME): ./cmd/api
	@echo ">>> Building $(APP_NAME)..."
	@go build -o ./bin/$(APP_NAME) ./cmd/api
	@echo ">>> SUCCESS..."

run: build myenv
	@echo ">>> Running $(APP_NAME)..."
	./bin/$(APP_NAME)

myenv:
	@echo "REQUIRED ENV"
	@echo "DB_URL: $(DB_URL)"

clean:
	@echo "CLEANING UP"
	rm -f ./bin/$(APP_NAME)
