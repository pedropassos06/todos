COMPOSE=docker compose

.PHONY: build backend-build frontend-build package test clean start restart stop logs ps

build: backend-build frontend-build

backend-build:
	$(MAKE) -C backend build

frontend-build:
	npm --prefix frontend run build

package:
	$(MAKE) -C backend package

test:
	$(MAKE) -C backend test

clean:
	$(MAKE) -C backend clean

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
