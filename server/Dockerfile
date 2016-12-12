FROM golang:1.7.4-alpine

RUN apk add --no-cache g++

ADD dockerize /usr/local/bin/dockerize

COPY vendor/src src/vendor
COPY src src
RUN go install ./src/...

ENTRYPOINT dockerize \
            -wait tcp://mqtt:1883 -timeout 30s \
            bin/shopping -com=/dev/tty.app -mqtt=mqtt:1883
