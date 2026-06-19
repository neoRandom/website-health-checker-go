FROM golang:1.24-alpine AS builder

WORKDIR /app

# Cache dependency downloads separately from source
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a statically linked binary (no CGO, no libc dependency)
RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-s -w" \
      -trimpath \
      -o /app/server \
      ./cmd/server

FROM scratch

# Pull TLS certificates so outbound HTTPS calls work
# This may required additional Prometheus configuration
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run as a non-root user (uid 65534 = nobody, minimal privileges)
USER 65534:65534

COPY --from=builder /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]