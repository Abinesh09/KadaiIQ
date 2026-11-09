# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api-server ./cmd/api/main.go

# Production scratch/alpine container
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/api-server /app/api-server

ENV PORT=8080
ENV APP_ENV=production

EXPOSE 8080

CMD ["/app/api-server"]

