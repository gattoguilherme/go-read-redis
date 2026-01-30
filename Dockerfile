# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /read-redis

# Final stage
FROM alpine:latest
COPY --from=builder /read-redis /read-redis
EXPOSE 3001
CMD ["/read-redis"]
