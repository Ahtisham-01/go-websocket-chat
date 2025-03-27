#  Golang image 
FROM golang:1.22.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache build-base

# Set working directory
WORKDIR /app

# Copy go mod
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application with verbose output
RUN CGO_ENABLED=0 GOOS=linux go build -v -o chat-server ./cmd/server

# Start a new stage from scratch
FROM alpine:latest

# Install necessary certificates and dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/chat-server .
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./chat-server"]