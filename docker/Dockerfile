# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /src

# Cache modules (copy only go.mod and go.sum)
COPY go.mod go.sum ./
RUN apk add --no-cache git ca-certificates && \
    go env -w GOPROXY=https://proxy.golang.org && \
    go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /app/kenjix ./cmd/api/main.go

# Final stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/kenjix /app/kenjix

# Use a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 8080
ENTRYPOINT ["/app/kenjix"]

