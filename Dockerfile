FROM golang:1.13.7-alpine AS go-builder

# Install build dependencies for docker-gen TODO:
RUN apk add --update \
        curl \
        gcc \
        git \
        make \
        musl-dev

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /app


FROM alpine:3.10

LABEL maintainer="Arthur Petukhovsky <petuhovskiy@yandex.ru> (@petuhovskiy)"

ENV DEBUG=false \
    DOCKER_HOST=unix:///var/run/docker.sock

# Install packages required by the image
RUN apk add --update \
        bash \
        ca-certificates \
        coreutils \
        curl \
        jq \
        openssl \
    && rm /var/cache/apk/*

COPY --from=go-builder /app ./

CMD [ "./app" ]
