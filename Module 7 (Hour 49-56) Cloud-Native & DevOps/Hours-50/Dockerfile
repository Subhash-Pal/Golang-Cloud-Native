FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

ARG VERSION=1.0.0
ARG COMMIT=local

RUN go build -ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT}" -o app .

FROM alpine:latest

WORKDIR /app

RUN adduser -D appuser

COPY --from=builder /app/app /app/app

EXPOSE 8080

USER appuser

CMD ["./app"]
