BIN_DIR=bin
BINARY=$(BIN_DIR)/bootstrap
ZIP_FILE=function.zip
COMPOSE=docker compose --project-directory .. -f ../docker-compose.yml

.PHONY: start restart stop logs ps build package test clean

start:
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env from .env.example"; fi
	$(COMPOSE) up --build -d --remove-orphans

restart:
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env from .env.example"; fi
	$(COMPOSE) down --remove-orphans
	$(COMPOSE) up --build -d --remove-orphans

stop:
	$(COMPOSE) down --remove-orphans

logs:
	$(COMPOSE) logs -f server

ps:
	$(COMPOSE) ps

build:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BINARY) ./cmd/lambda

package: build
	zip -j $(ZIP_FILE) $(BINARY)

test:
	go test ./...

clean:
	rm -f $(BINARY) $(ZIP_FILE)
	rmdir --ignore-fail-on-non-empty $(BIN_DIR)
