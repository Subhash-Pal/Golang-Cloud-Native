FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app /app/app

EXPOSE 8080

CMD ["./app"]
