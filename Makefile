BIN_DIR=bin
BINARY=$(BIN_DIR)/bootstrap
ZIP_FILE=function.zip
COMPOSE=docker compose

.PHONY: build package clean start restart stop logs ps

build:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BINARY) ./cmd/lambda

package: build
	zip -j $(ZIP_FILE) $(BINARY)

clean:
	rm -f $(BINARY) $(ZIP_FILE)
	rmdir --ignore-fail-on-non-empty $(BIN_DIR)

start:
	$(COMPOSE) up --build -d --remove-orphans

restart:
	$(COMPOSE) down --remove-orphans
	$(COMPOSE) up --build -d --remove-orphans

stop:
	$(COMPOSE) down --remove-orphans

logs:
	$(COMPOSE) logs -f server

ps:
	$(COMPOSE) ps
