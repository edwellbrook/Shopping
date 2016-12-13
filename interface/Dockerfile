FROM golang:1.7.4-alpine

RUN apk add --no-cache g++

COPY vendor/src src/vendor
COPY src src
RUN go install ./src/...

ENTRYPOINT bin/serial-receiver \
            -com=/dev/tty.app \
            -mqtt=rabbitmq:1883 \
            -psql=postgres://postgres@postgres:5432?sslmode=disable
