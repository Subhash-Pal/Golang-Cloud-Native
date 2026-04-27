FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o api .

EXPOSE 8080

CMD ["./api"]
