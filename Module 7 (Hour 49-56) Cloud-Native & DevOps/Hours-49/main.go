package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type message struct {
	Hour        int      `json:"hour"`
	Topic       string   `json:"topic"`
	Description string   `json:"description"`
	Endpoints   []string `json:"endpoints"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, message{
			Hour:        49,
			Topic:       "Docker fundamentals",
			Description: "A simple Go HTTP API that is easy to package into a Docker image.",
			Endpoints:   []string{"/", "/healthz", "/time"},
		})
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"server_time": time.Now().Format(time.RFC3339)})
	})

	addr := ":" + port
	log.Printf("hour 49 server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

/*
Invoke the server with:
$ go run main.go

Next Console enter below command to set the PORT environment variable:
$ export PORT=8080
Then, you can test the endpoints with:
$ curl http://localhost:8080/
$ curl http://localhost:8080/healthz
$ curl http://localhost:8080/time
*/
//build the Docker image with:
//$ docker build -t hour49:latest .
//Run the container with:
//$ docker run -p 8080:8080 hour49:latest	

/*
To set up a Go application with PostgreSQL and ensure your data is preserved (persisted), you need a multi-stage Dockerfile and a Docker Compose file that uses volumes.
1. Dockerfile (for your Go App)
This multi-stage Dockerfile builds your application in a Go environment and then copies only the binary and migration files into a lightweight production image.
dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Run stage
FROM alpine:latest
WORKDIR /root/
# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .
# Copy your SQL migration files so the app can run them on startup
COPY ./migrations ./migrations

CMD ["./main"]
Use code with caution.
2. Docker Compose (for the Stack)
The Docker Compose file defines your services and uses a named volume to ensure PostgreSQL data is not lost when containers are stopped or deleted.
yaml
services:
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d myapp"]
      interval: 5s
      timeout: 5s
      retries: 5
    volumes:
      # PERSISTENCE: This saves your Postgres data to a Docker volume
      - postgres_data:/var/lib/postgresql/data
      # INITIAL SEED: (Optional) Any .sql file here runs only on the VERY FIRST start
      - ./init-db:/docker-entrypoint-initdb.d

  app:
    build: .
    depends_on:
      db:
        condition: service_healthy
    environment:
      # The app uses the service name 'db' as the hostname
      DATABASE_URL: postgres://user:password@db:5432/myapp?sslmode=disable

volumes:
  postgres_data:
Use code with caution.
3. Copying Existing PostgreSQL Data
If you already have a PostgreSQL data folder on your machine and want to "copy" it into Docker:
Bind Mounts: Instead of postgres_data:/var/lib/postgresql/data, use a direct path: - ./my-local-data:/var/lib/postgresql/data. This mounts your local folder directly into the container.
Database Dumps: If you have a .sql backup, place it in the ./init-db folder mentioned above. When the database container starts for the first time, it will automatically execute that script to recreate your tables and data.
Manual Import: You can copy a file into a running container using docker cp backup.sql db:/backup.sql and then run docker exec -it db psql -U user -d myapp -f /backup.sql.

*/