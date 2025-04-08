# Use the official Golang image as the base image
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN GOOS=linux go build -o breviago

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/breviago .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./breviago"]
