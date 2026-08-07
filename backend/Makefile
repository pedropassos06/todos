BIN_DIR=bin
BINARY=$(BIN_DIR)/bootstrap
ZIP_FILE=$(BIN_DIR)/function.zip
AWS_REGION?=us-east-1
FUNCTION_NAME?=todos-api
COMPOSE=docker compose --project-directory .. -f ../docker-compose.yml

.PHONY: start restart stop logs ps build package deploy-lambda test clean

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

deploy-lambda: package
	aws lambda update-function-code --function-name "$(FUNCTION_NAME)" --zip-file "fileb://$(ZIP_FILE)" --region "$(AWS_REGION)"
	aws lambda wait function-updated --function-name "$(FUNCTION_NAME)" --region "$(AWS_REGION)"

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR) $(ZIP_FILE)
