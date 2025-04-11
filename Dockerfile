# syntax=docker/dockerfile:1

ARG GOVERSION=1.24.2
ARG PKL_VERSION=0.28.1

FROM golang:${GOVERSION} AS dev
RUN go install "github.com/air-verse/air@latest"
RUN go install "github.com/pressly/goose/v3/cmd/goose@latest"
RUN go install "github.com/apple/pkl-go/cmd/pkl-gen-go@latest"
RUN go install "github.com/a-h/templ/cmd/templ@latest"
ARG PKL_VERSION
RUN curl -L -o /usr/bin/pkl "https://github.com/apple/pkl/releases/download/${PKL_VERSION}/pkl-linux-$(uname -m)" && chmod +x /usr/bin/pkl
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ --mount=type=bind,target=.
CMD ["air", "-c", ".air.toml"]


FROM golang:${GOVERSION}-alpine AS prod
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ARG TARGETARCH
ADD --chmod=111 "https://github.com/apple/pkl/releases/download/${PKL_VERSION}/pkl-alpine-linux-${TARGETARCH}" /bin/pkl
ARG SOURCE_COMMIT
RUN echo $SOURCE_COMMIT > release.txt
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /bin/manage ./cmd/manage/
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server/
COPY . .
CMD ["/bin/server"]
