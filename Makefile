SOURCE_COMMIT_SHA := $(shell git rev-parse HEAD)

ENVS := SOURCE_COMMIT=${SOURCE_COMMIT_SHA} COMPOSE_BAKE=true


.PHONY: run dev build-dev fdev prod fprod logs-prod go-to-server-container pkl-gen db-status db-up db-down db-reset templ add-pill-dispenser build-prod-image

run: dev

dev:
	${ENVS} docker compose -f compose.yaml up

build-dev:
	${ENVS} docker compose -f compose.yaml up --build

fdev:
	${ENVS} docker compose -f compose.yaml down

prod:
	docker compose -f compose.test-prod.yaml up --build -d

fprod:
	docker compose -f compose.test-prod.yaml down

logs-prod:
	docker compose -f compose.test-prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it --tty pill-dispenser-agent /bin/sh

db-status:
	docker exec -it --tty pill-dispenser-agent goose postgres "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations status

db-up:
	docker exec -it --tty pill-dispenser-agent goose postgres "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations up

db-down:
	docker exec -it --tty pill-dispenser-agent goose postgres "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations down

db-reset:
	docker exec -it --tty pill-dispenser-agent goose postgres "$(shell docker exec -it pill-dispenser-agent /bin/manage -c print-db-string)" -dir=internal/db/migrations reset

pkl-gen:
	docker exec -it --tty pill-dispenser-agent pkl-gen-go pkl/config.pkl --base-path github.com/tikhonp/medsenger-pill-dispenser-bot

templ:
	docker exec -it --tty pill-dispenser-agent templ generate

add-pill-dispenser:
	docker exec -it --tty pill-dispenser-agent manage -c add-pill-dispenser -i

build-prod-image:
	docker buildx build --build-arg SOURCE_COMMIT="${SOURCE_COMMIT_SHA}" --target prod -t docker.telepat.online/agents-pilldispenser-image:latest .

update-deps:
	docker exec -it --tty pill-dispenser-agent go get -u ./...
