# syntax=docker/dockerfile:1

ARG GOVERSION=1.24.6

FROM golang:${GOVERSION}-alpine AS dev
RUN go install "github.com/air-verse/air@latest" && \
    go install "github.com/a-h/templ/cmd/templ@latest" &&\
    go install "github.com/pressly/goose/v3/cmd/goose@latest"
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
CMD ["air", "-c", ".air.toml"]


FROM golang:${GOVERSION}-alpine AS build-prod
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH \
    go install -tags='no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' github.com/pressly/goose/v3/cmd/goose@latest
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/manage ./cmd/manage/
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server/

FROM alpine AS prod
WORKDIR /src
COPY --from=build-prod /bin/server /bin/manage /go/bin/goose /bin/
COPY . .
EXPOSE 80
ENV DEBUG=false
ARG SOURCE_COMMIT
ENV SOURCE_COMMIT=${SOURCE_COMMIT}
ENV SERVER_PORT=80
ENTRYPOINT ["/bin/sh", "-c", "goose postgres \"$(manage -c print-db-string)\" -dir=internal/db/migrations up && server"]
