# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add tzdata

# Copy the binary and env file from builder
COPY --from=builder /build/main /main
COPY .env /.env

# Expose port 3000
EXPOSE 3000

# Run the application
CMD ["/main"]
