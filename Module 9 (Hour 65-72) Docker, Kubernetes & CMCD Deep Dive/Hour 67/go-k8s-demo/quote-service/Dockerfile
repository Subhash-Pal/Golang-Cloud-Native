FROM golang:1.26-alpine
WORKDIR /app
COPY . .
RUN go build -o quote-service .
EXPOSE 8081
CMD ["./quote-service"]