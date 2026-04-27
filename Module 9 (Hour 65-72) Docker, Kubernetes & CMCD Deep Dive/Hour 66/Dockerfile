FROM golang:1.26-alpine AS builder

WORKDIR /src

# Cache module metadata separately to speed up rebuilds.
COPY go.mod ./
RUN go mod download

COPY main.go ./

# Build a small static binary for a minimal runtime image.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app ./main.go
# what is -trimpath -ldflags="-s -w" in go build command?
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# Runtime security defaults.
ENV PORT=8080

COPY --from=builder /out/app /app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app"]
