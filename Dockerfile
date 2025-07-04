# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 is important for creating a statically linked binary
# -ldflags="-s -w" reduces the binary size by removing debug information
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o secusend-backend .

# Stage 2: Create the final, minimal image
FROM scratch

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/secusend-backend .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose the port your Fiber app listens on (default is 3000)
EXPOSE 3000

# Command to run the executable
CMD ["./secusend-backend"]
