# =========================================================================
# Stage 1: Builder
# =========================================================================
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git curl libstdc++ libgcc

WORKDIR /app

# Install Go tools
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install Tailwind CSS standalone binary (musl-compatible for Alpine)
ARG TAILWIND_VERSION=v4.1.18
ARG TARGETARCH
RUN set -eux; \
    ARCH="${TARGETARCH:-$(apk --print-arch)}"; \
    case "${ARCH}" in \
    amd64|x86_64) TW_ARCH="x64" ;; \
    arm64|aarch64) TW_ARCH="arm64" ;; \
    *) echo "Unsupported architecture: ${ARCH}" >&2; exit 1 ;; \
    esac; \
    TW_MUSL_URL="https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/tailwindcss-linux-${TW_ARCH}-musl"; \
    TW_GLIBC_URL="https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/tailwindcss-linux-${TW_ARCH}"; \
    if ! curl -fsSL "${TW_MUSL_URL}" -o /usr/local/bin/tailwindcss; then \
    curl -fsSL "${TW_GLIBC_URL}" -o /usr/local/bin/tailwindcss; \
    fi; \
    chmod +x /usr/local/bin/tailwindcss

# Generate Templ templates
RUN templ generate

# Build CSS for production
RUN tailwindcss -i assets/css/index.css -o assets/css/output.css --minify

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# =========================================================================
# Stage 2: Runner
# =========================================================================
FROM alpine:3.19 AS runner

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy binary and assets from builder
COPY --from=builder /app/server .
# Copy only the compiled CSS and static assets
COPY --from=builder /app/assets/css/output.css ./assets/css/output.css
COPY --from=builder /app/assets/static ./assets/static

# Set ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./server"]
