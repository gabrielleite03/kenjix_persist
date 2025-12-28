# Build stage
FROM golang:1.24.5-alpine AS builder
WORKDIR /src

# Cache modules (copy only go.mod so build context works when go.sum is absent)
COPY go.mod ./
RUN apk add --no-cache git ca-certificates && \
	go env -w GOPROXY=https://proxy.golang.org && \
	go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /app/kenjix ./cmd/api

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

