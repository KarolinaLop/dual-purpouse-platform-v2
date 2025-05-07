# -------- Build Stage --------
    FROM golang:1.24 AS builder

    ENV CGO_ENABLED=1
    
    # Install dependencies needed for cgo and sqlite3
    RUN apt-get update && apt-get install -y gcc libc6-dev
    
    WORKDIR /app
    
    # Copy go mod files and download dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY app.db /app/app.db
    
    # Copy the rest of the app
    COPY . .
    
    # Build the Go app
    RUN go build -v -o /app/app ./cmd/DP/main.go
    
    
    # -------- Runtime Stage --------
    FROM debian:bookworm-slim
    
    RUN apt-get update && apt-get install -y nmap libsqlite3-0 && apt-get clean
    
    COPY --from=builder /app/app /usr/local/bin/app
    COPY --from=builder /app/app.db app.db
    EXPOSE 8080
    
    ENV GIN_MODE=release
    
    CMD ["app"]