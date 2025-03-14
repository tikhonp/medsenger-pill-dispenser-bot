run: dev

dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
dev:
	docker compose -f compose.yaml up

build-dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
build-dev:
	docker compose -f compose.yaml up --build

prod: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
prod:
	docker compose -f compose.yaml up --build -d

fprod:
	docker compose -f compose.yaml down

logs-prod:
	docker compose -f compose.yaml logs -f -n 100

go-to-server-container:
	docker exec -it pill-dispenser-server-agent /bin/bash
