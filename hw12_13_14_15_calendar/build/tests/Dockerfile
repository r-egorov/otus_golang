FROM golang:1.18-alpine

WORKDIR /integration_tests

COPY ./integration_tests /integration_tests

RUN apk update && apk add git build-base
RUN go install github.com/cucumber/godog/cmd/godog@latest

ENTRYPOINT ["godog", "run"]