# Hour 49 keeps Docker simple: one image that both builds and runs the app.
FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]
