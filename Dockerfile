# Build stage
FROM golang:1.24.12-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/config ./internal/config

EXPOSE 8080

CMD ["./main"]