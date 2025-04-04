.PHONY: run dev build-dev prod fprod logs-prod go-to-server-container pkl-gen db-status db-up db-down db-reset templ

run: dev

dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
dev:
	docker compose -f compose.yaml up

build-dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD) 
build-dev: export COMPOSE_BAKE=true
build-dev:
	docker compose -f compose.yaml up --build

prod: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
prod:
	docker compose -f compose.prod.yaml up --build -d

fprod:
	docker compose -f compose.prod.yaml down

logs-prod:
	docker compose -f compose.prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it pill-dispenser-server-agent /bin/bash

db-status:
	docker exec -it pill-dispenser-agent goose sqlite3 "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations status

db-up:
	docker exec -it pill-dispenser-agent goose sqlite3 "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations up

db-down:
	docker exec -it pill-dispenser-agent goose sqlite3 "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations down

db-reset:
	docker exec -it pill-dispenser-agent goose sqlite3 "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations reset

pkl-gen:
	pkl-gen-go pkl/config.pkl --base-path github.com/tikhonp/medsenger-pill-dispenser-bot

templ:
	templ generate
