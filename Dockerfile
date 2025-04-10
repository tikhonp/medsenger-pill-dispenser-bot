# syntax=docker/dockerfile:1

ARG GOVERSION=1.24.2

FROM golang:${GOVERSION}-bookworm AS dev
WORKDIR /src
ADD --chmod=111 "https://github.com/apple/pkl/releases/download/0.28.1/pkl-linux-aarch64" /usr/bin/pkl
RUN go install "github.com/air-verse/air@latest"
RUN go install "github.com/pressly/goose/v3/cmd/goose@latest"
RUN go install "github.com/apple/pkl-go/cmd/pkl-gen-go@latest"
RUN go install "github.com/a-h/templ/cmd/templ@latest"
COPY go.mod go.sum ./
RUN go mod download && go mod verify
ARG SOURCE_COMMIT
RUN echo $SOURCE_COMMIT > /src/release.txt
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /bin/manage ./cmd/manage/
COPY . .
CMD ["air"]


FROM golang:${GOVERSION} AS prod
WORKDIR /src
ADD --chmod=111 'https://github.com/apple/pkl/releases/download/0.28.1/pkl-alpine-linux-amd64' /bin/pkl
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ARG SOURCE_COMMIT
RUN echo $SOURCE_COMMIT > release.txt
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=.\
    GOARCH=$TARGETARCH go build -o /bin/manage ./cmd/manage/
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server/
COPY . .
CMD ["/bin/server"]
