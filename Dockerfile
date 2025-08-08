# syntax=docker/dockerfile:1
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY .env /app/cmd/.env

RUN go mod download

COPY . .
COPY .env /app/cmd/.env

RUN go build -o event-pipeline cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/event-pipeline .
EXPOSE 9000

CMD ["./event-pipeline"]