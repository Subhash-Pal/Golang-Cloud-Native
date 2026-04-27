# Use the official Golang image with version 1.26 or higher
FROM golang:1.26-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application code
COPY . .

# Build the Go application
RUN go build -o app .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./app"]