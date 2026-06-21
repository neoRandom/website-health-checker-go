FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a statically linked binary (no CGO, no libc dependency)
RUN GOMEMLIMIT=200MiB CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-s -w" \
      -trimpath \
      -o /app/server \
      ./cmd/server

FROM alpine:3.22

RUN apk add --no-cache su-exec

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /app/server
COPY --from=builder /app/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

# Fallback if no volume is attached
RUN mkdir -p /app/data && chown -R 65534:65534 /app/data

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]