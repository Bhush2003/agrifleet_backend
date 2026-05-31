# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Cache dependencies layer
COPY go.mod go.sum ./
RUN go mod download

# Build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o agrifleet ./cmd/server

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/agrifleet .

# Railway injects PORT env var — our app reads APP_PORT
EXPOSE 8080

CMD ["./agrifleet"]
