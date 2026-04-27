# Stage 1: Build the Go application
FROM golang:1.26-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go application code
COPY . .

# Build the Go application
RUN go build -o app .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./app"]