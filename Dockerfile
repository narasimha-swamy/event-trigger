# Build stage
FROM golang:1.22-alpine AS builder

# Install required tools
RUN apk add --no-cache git

# Install specific swag version
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Generate Swagger docs
RUN swag init -g main.go --output docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy built binary and assets
COPY --from=builder /app/backend .
COPY --from=builder /app/web ./web
COPY --from=builder /app/docs ./docs

# Environment variables
ENV DATABASE_URL="host=db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
ENV PORT=8080

EXPOSE 8080

CMD ["./backend"]