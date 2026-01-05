# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Build the import-cities tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o import-cities ./cmd/import-cities

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/import-cities .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]

