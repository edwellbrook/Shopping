FROM golang:1.7.4-alpine

RUN apk add --no-cache g++

COPY vendor/src src/vendor
COPY src src
RUN go install ./src/...

ENV COM /dev/tty.app
ENV MQTT 127.0.0.1:1883
ENV PSQL postgres://postgres@127.0.0.1:5432

ENTRYPOINT bin/serial-receiver \
            -com=$COM \
            -mqtt=$MQTT \
            -psql=$PSQL
