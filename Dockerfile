FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

RUN apt-get update && \
    apt-get install -y libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 8080

CMD ["./main"]
