# STAGE 1: Build the binary
FROM golang:1.26-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker caching
COPY go.mod ./
# If you have a go.sum, copy it too: COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# STAGE 2: Create the final minimal image
FROM alpine:latest  
# Or use 'scratch' for an even smaller, empty base image

WORKDIR /root/

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
