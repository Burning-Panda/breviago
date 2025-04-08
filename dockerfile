# TODO: https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

# Use the official Golang image as the base image
FROM golang:1.24-bookworm as builder
ENV CGO_ENABLED=1

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

# Install necessary runtime dependencies
RUN apk add --no-cache libc6-compat

# Set the working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/breviago .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./breviago"]
