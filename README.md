# Medsenger Pill Dispenser Agent

Agent for managing telepat smart pill dispenser scheduling.

# Development

1. Install __docker__ and __make__

2. Create `config.pkl` file in project root folder amendig `pkl/congig.pkl` scheme

### Run Development

```sh
make
```

or

```sh
make dev
```

or 

```sh
make build-dev # preferred if config files were changed, so it rebuilds image
```

### Apply migrations

```sh
make db-status # check
make db-up     # apply all migrations
```

I use [goose](https://github.com/pressly/goose) to manage migrations. To create new migration run:

```sh
goose -dir=internal/db/migrations create <migration-name> sql
```

Other db shortcuts:

```sh
make db-reset # reset all migrations (gooose reset)
make db-down  # goose down
```

### Modifying config scheme

Config managed by tool [pkl](https://pkl-lang.org/index.html)

If you want to modify config scheme edit file/files in `pkl` directory. After regenerate go code:

```sh
make pkl-gen
```

> development docker container must be active

### HTML templating

I use [templ](https://github.com/a-h/templ) as template engine. After changing `*.templ` files regenerate go code using:

```sh
make templ
```

> development docker container must be active

### Enter server container shell

There is shortcut for this:

```sh
make go-to-server-container
```

# Deploying

To deploy you also need __docker__ and __make__. In project root run:

```sh
make prod
```

It will create prod containers and run it in detached mode.

To stop run:

```sh
make fprod
```

To view logs in real time:

```sh
make logs-prod
```

# License

Created by Tikhon Petrishchev 2025

Copyright © 2025 OOO Telepat. All rights reserved.
