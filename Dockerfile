# =========================================================================
# Stage 1: Builder
# =========================================================================
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git curl bash

WORKDIR /app

# Install Go tools
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install CSS Framework & Assets
RUN cd assets && curl -sL daisyui.com/fast | bash
RUN mkdir -p assets/css assets/js && \
    mv assets/tailwindcss .
RUN curl -sL -o assets/css/basecoat.min.css https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.css
RUN echo '@import "tailwindcss"; @import "./basecoat.min.css";' > assets/css/index.css
RUN rm assets/input.css assets/output.css assets/daisyui* 2>/dev/null || true
# Download JS
RUN curl -sL -o assets/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js && \
    curl -sL -o assets/js/hyperscript.min.js https://unpkg.com/hyperscript.org@0.9.14 && \
    curl -sL -o assets/js/basecoat.min.js https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.js

# Generate Templ templates
RUN templ generate

# Build CSS for production
RUN ./tailwindcss -i assets/css/index.css -o assets/css/output.css --minify

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
