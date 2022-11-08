FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o bin/migrator ./cmd/migrator

FROM alpine:3.16

COPY --from=builder /app/bin/migrator /app/bin/migrator

ENTRYPOINT ["/app/bin/migrator"]
