FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o bin/calendar ./cmd/calendar

FROM alpine:3.16

COPY --from=builder /app/bin/calendar /app/bin/calendar

EXPOSE 8080

ENTRYPOINT ["/app/bin/calendar"]
