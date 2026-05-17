FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/sso ./cmd/sso/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/migrator ./cmd/migrator/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata postgresql-client

COPY --from=builder /app/bin/sso /app/sso
COPY --from=builder /app/bin/migrator /app/migrator
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

RUN chmod +x /app/sso /app/migrator

EXPOSE 8081

CMD ["/app/sso"]