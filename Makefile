SOURCE_COMMIT_SHA := $(shell git rev-parse HEAD)

ENVS := SOURCE_COMMIT=${SOURCE_COMMIT_SHA} COMPOSE_BAKE=true


.PHONY: run dev build-dev prod fprod logs-prod go-to-server-container pkl-gen db-status db-up db-down db-reset templ add-pill-dispenser

run: dev

dev:
	${ENVS} docker compose -f compose.yaml up

build-dev:
	${ENVS} docker compose -f compose.yaml up --build

prod:
	${ENVS} docker compose -f compose.prod.yaml up --build -d

fprod:
	${ENVS} docker compose -f compose.prod.yaml down

logs-prod:
	${ENVS} docker compose -f compose.prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it --tty pill-dispenser-agent /bin/bash

db-status:
	docker exec -it --tty pill-dispenser-agent goose sqlite3 "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations status

db-up:
	docker exec -it --tty pill-dispenser-agent goose sqlite3 "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations up

db-down:
	docker exec -it --tty pill-dispenser-agent goose sqlite3 "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations down

db-reset:
	docker exec -it --tty pill-dispenser-agent goose sqlite3 "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations reset

pkl-gen:
	docker exec -it --tty pill-dispenser-agent pkl-gen-go pkl/config.pkl --base-path github.com/tikhonp/medsenger-pill-dispenser-bot

templ:
	docker exec -it --tty pill-dispenser-agent templ generate

add-pill-dispenser:
	docker exec -it --tty pill-dispenser-agent manage -c add-pill-dispenser -i
