# syntax=docker.io/docker/dockerfile-upstream:1.17.0-labs

ARG GOVERSION=1.25.0
ARG NODE_VERSION=24


FROM --platform=$BUILDPLATFORM node:${NODE_VERSION}-alpine AS frontend-builder
WORKDIR /app
ARG VITE_AGENT_URL
ENV VITE_AGENT_URL=${VITE_AGENT_URL}
COPY frontend/package.json ./
RUN --mount=type=cache,target=/root/.npm npm install
COPY frontend/ ./
RUN npm run build


FROM golang:${GOVERSION}-alpine AS dev
RUN go install "github.com/air-verse/air@latest" && \
    go install "github.com/a-h/templ/cmd/templ@latest" &&\
    go install "github.com/pressly/goose/v3/cmd/goose@latest"
COPY --from=frontend-builder /app/dist /public/frontend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
CMD ["air", "-c", ".air.toml"]


FROM --platform=$BUILDPLATFORM golang:${GOVERSION}-alpine AS build-prod
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -tags='no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' -o /bin/goose github.com/pressly/goose/v3/cmd/goose
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/manage ./cmd/manage/
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/server ./cmd/server/

FROM alpine AS prod
WORKDIR /src
COPY --from=build-prod /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
COPY --from=build-prod /bin/server /bin/manage /bin/goose /bin/
COPY --exclude=frontend/ . .
EXPOSE 80
ENV DEBUG=false
ARG SOURCE_COMMIT
ENV SOURCE_COMMIT=${SOURCE_COMMIT}
ENV SERVER_PORT=80
ENTRYPOINT ["/bin/sh", "-c", "goose postgres \"$(manage -c print-db-string)\" -dir=internal/db/migrations up && server"]
