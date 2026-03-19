# Use the official Golang image as a base image
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Cache modules (copy only go.mod and go.sum)
COPY go.mod go.sum ./
RUN apk add --no-cache git ca-certificates && \
    go env -w GOPROXY=https://proxy.golang.org && \
    go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/api

# Use a minimal base image for the final container
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]

