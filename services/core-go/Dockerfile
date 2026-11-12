# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and build essentials
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build single static executable binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o core_server ./cmd/server/main.go

# Production minimal scratch/alpine stage
FROM alpine:3.19

WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/core_server .

EXPOSE 8080

CMD ["./core_server"]
