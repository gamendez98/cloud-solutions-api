
# Use official golang image as a builder
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -a -installsuffix cgo -o main .

# Use a minimal image for the final artifact
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*


# Set working directory for the app
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .
COPY docker.env .env
COPY api-service-account.json api-service-account.json

RUN mkdir "uploads"

# Expose application port
EXPOSE 80
EXPOSE 8080

# Command to run the executable
CMD ["./main"]