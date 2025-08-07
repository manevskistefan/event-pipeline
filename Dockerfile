# syntax=docker/dockerfile:1
FROM golang:1.24.5-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
RUN go build -o event-pipeline cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/event-pipeline .
EXPOSE 8080

CMD ["./event-pipeline"]