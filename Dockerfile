# Stage 1: Build the application binary safely
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install system dependencies needed for compiling Go toolchains if necessary
RUN apk add --no-cache git gcc musl-dev

# Leverage Docker cache layer for module requirements
COPY go.mod go.sum ./
RUN go mod download

# Copy the remaining project layout sources
COPY . .

# Compile target binary directly without symbols for reduced container sizes
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/server cmd/server/main.go

# Stage 2: Distribute light execution runtime image
FROM alpine:3.18

WORKDIR /app

# Import compiled binaries out of the builder image
COPY --from=builder /app/bin/server .

# Expose standard REST traffic pipeline port
EXPOSE 8080

ENTRYPOINT ["./server"]
