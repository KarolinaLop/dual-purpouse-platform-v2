FROM alpine:latest

# Install Nmap
RUN apk add --no-cache nmap

# Copy the application
WORKDIR /app
COPY . .

# Run the application
CMD ["go", "run", "main.go"]