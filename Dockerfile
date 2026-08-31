# Stage 1: Build the Go application, means go and download this
FROM golang:1.26.4 AS builder

WORKDIR /app

# Copy dependency files and vendored dependencies
COPY go.mod go.sum ./
COPY vendor ./vendor

# Copy application source code
COPY . .

# Build the application
RUN go build -mod=vendor -o rssaggregator .

# Stage 2: Run the application
FROM debian:bookworm-slim

WORKDIR /app

# Copy only the compiled binary from the builder
COPY --from=builder /app/rssaggregator .

EXPOSE 8080

CMD ["./rssaggregator"]