# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Install air
RUN go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Final stage
FROM golang:1.24.1-alpine

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache make git gcc musl-dev

# Copy air binary from builder
COPY --from=builder /go/bin/air /go/bin/air

# Copy the source code and dependencies
COPY --from=builder /go/pkg /go/pkg
COPY --from=builder /app .

COPY internal/ledger/infrastructure/database/payments.json /app/internal/ledger/infrastructure/database/
COPY internal/ledger/infrastructure/database/expenses.json /app/internal/ledger/infrastructure/database/

# Copy .env file
COPY .env /.env

# Expose port 3000
EXPOSE 3000

# Run the application with the correct main path
CMD ["air", "-c", ".air.toml"]

