# Use the golang image to build the application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Use a minimal image to run the application
FROM alpine:3.15

WORKDIR /app

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./main"]
