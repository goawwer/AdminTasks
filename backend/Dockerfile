FROM golang:1.24.2 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN wget -O /app/migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \
    && tar -xzf /app/migrate.tar.gz -C /app \
    && rm /app/migrate.tar.gz

RUN go build -v -o app ./cmd/AdminTasksAPI

FROM debian:stable-slim
WORKDIR /app

RUN apt-get update && apt-get install -y postgresql-client --no-install-recommends && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app .
COPY --from=builder /app/migrate .
COPY migrations ./migrations
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x /app/entrypoint.sh /app/migrate

EXPOSE 8080

CMD ["./entrypoint.sh"]