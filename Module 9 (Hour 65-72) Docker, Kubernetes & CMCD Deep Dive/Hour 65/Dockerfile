# Stage 1: Build the Go application
FROM golang:1.26-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go application code
COPY . .

# Build the Go application
RUN go build -o app .

# Stage 2: Create a minimal runtime image
FROM alpine:3.18

# Create a non-root user
RUN adduser -D appuser

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Switch to the non-root user
USER appuser

# Expose the port
EXPOSE 8080

# Add a health check
HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget -q -O - http://localhost:8080 || exit 1

# Run the application
CMD ["./app"]