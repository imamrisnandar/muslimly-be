# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies (cached)
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Stage 2: Runtime (Minimal)
FROM alpine:latest

WORKDIR /root/

# Install CA Certs for Firebase/External APIs
RUN apk --no-cache add ca-certificates tzdata

# Copy Binary and Configs
COPY --from=builder /app/main .
# COPY --from=builder /app/firebase-service-account.json . # Don't copy this if using volume mount
# COPY --from=builder /app/config.yaml . # Don't copy secrets!

# Expose Port
EXPOSE 8080

# Run
CMD ["./main"]
