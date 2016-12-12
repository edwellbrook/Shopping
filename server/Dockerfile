FROM golang:1.7.4-alpine

RUN apk add --no-cache g++

COPY vendor/src src/vendor
COPY src src
RUN go install ./src/...

ENTRYPOINT bin/shopping -com=/dev/tty.app -mqtt=mqtt:1883
