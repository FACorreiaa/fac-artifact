# GoForge Project Makefile
# =========================================================================
# Usage:
#   make setup    - Install all tools (Air, Templ, Goose, Tailwind + CSS Framework)
#   make dev      - Start development server with live reload
#   make build    - Build production binary
#   make test     - Run tests
# =========================================================================

PROJECT_NAME := myapp
BINARY_NAME := server



.PHONY: all build run test clean dev setup help

all: build

# =========================================================================
# Setup & Dependencies
# =========================================================================

setup: ## Install all development tools
	@echo "📦 Installing Go tools..."
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest

	@echo "🧹 Tidying modules..."
	@go mod tidy
	@echo "📥 Installing Tailwind CSS + Basecoat..."
	@mkdir -p assets/css assets/js
	@cd assets && curl -sL daisyui.com/fast | bash
	@mv assets/tailwindcss tailwindcss 2>/dev/null || true
	@echo "Downloading Basecoat CSS..."
	@curl -sL -o assets/css/basecoat.min.css https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.css
	@echo "Creating assets/css/index.css..."
	@echo '@import "tailwindcss"; @import "./basecoat.min.css";' > assets/css/index.css
	@echo "Cleaning up DaisyUI files..."
	@rm assets/input.css assets/output.css assets/daisyui.mjs assets/daisyui-theme.mjs 2>/dev/null || true
	@echo "📥 Downloading Frontend Libraries..."
	@curl -sL -o assets/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js
	@curl -sL -o assets/js/hyperscript.min.js https://unpkg.com/hyperscript.org@0.9.14
	@curl -sL -o assets/js/basecoat.min.js https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.js
	@echo ""
	@echo "✅ Setup complete! Run 'make dev' to start development."

tidy: ## Tidy Go modules
	go mod tidy

ci-setup: ## Setup for CI environments
	@echo "📦 Installing Go tools for CI..."
	go install github.com/a-h/templ/cmd/templ@latest
	@echo "🧹 Tidying modules..."
	@go mod tidy
	@echo "📥 Installing Tailwind CSS + Basecoat..."
	@mkdir -p assets/css assets/js
	@cd assets && curl -sL daisyui.com/fast | bash
	@mv assets/tailwindcss tailwindcss 2>/dev/null || true
	@echo "Downloading Basecoat CSS..."
	@curl -sL -o assets/css/basecoat.min.css https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.css
	@echo "Creating assets/css/index.css..."
	@echo '@import "tailwindcss"; @import "./basecoat.min.css";' > assets/css/index.css
	@echo "Cleaning up DaisyUI files..."
	@rm assets/input.css assets/output.css assets/daisyui.mjs assets/daisyui-theme.mjs 2>/dev/null || true
	@echo "📥 Downloading Frontend Libraries..."
	@curl -sL -o assets/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js
	@curl -sL -o assets/js/hyperscript.min.js https://unpkg.com/hyperscript.org@0.9.14
	@curl -sL -o assets/js/basecoat.min.js https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.js
	@echo "✅ CI setup complete!"

# =========================================================================
# Development
# =========================================================================

dev: ## Start development server with live reload (using Air)
	@echo "🚀 Starting development server with Air..."
	@GO_ENV=development air

dev-templ: ## Start development with Templ proxy (auto browser refresh)
	@echo "🚀 Starting with Templ proxy..."
	@echo "📡 Access via: http://127.0.0.1:7331 (auto-refresh enabled)"
	@echo "📡 Direct access: http://localhost:8080 (manual refresh needed)"
	@GO_ENV=development make -j2 dev-templ-watch dev-tailwind-watch

dev-templ-watch:
	@templ generate --watch --proxy="http://localhost:8080" --cmd="go run ./cmd/server"

dev-tailwind-watch:
	@sleep 2 && tailwindcss -i ./assets/css/index.css -o ./assets/css/output.css --watch

templ: ## Generate Templ templates once
	templ generate

css: ## Build CSS once
	@tailwindcss -i assets/css/index.css -o assets/css/output.css

# =========================================================================
# Build
# =========================================================================

build: templ ## Build production binary
	@echo "🔨 Building CSS..."
	@tailwindcss -i assets/css/index.css -o assets/css/output.css --minify
	@echo "🔨 Building binary..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/$(BINARY_NAME) ./cmd/server
	@echo "✅ Build complete: ./bin/$(BINARY_NAME)"

run: build ## Build and run the application
	./bin/$(BINARY_NAME)

clean: ## Remove build artifacts
	rm -rf bin/
	rm -rf tmp/
	rm -f assets/css/output.css
	rm -f tailwindcss
	find . -name "*_templ.go" -delete

# =========================================================================
# Testing
# =========================================================================

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run golangci-lint
	golangci-lint run



# =========================================================================
# Deployment
# =========================================================================

deploy: ## Deploy to production server (Hetzner VPS)
	@echo "🚀 Deploying to Hetzner VPS..."
	@./deploy/deploy.sh

deploy-fly: ## Deploy to Fly.io
	@echo "🚀 Deploying to Fly.io..."
	@fly deploy

deploy-setup: ## First-time server setup instructions
	@echo "📖 Deployment options:"
	@echo ""
	@echo "   🆓 Fly.io (FREE): deploy/FLY-IO.md"
	@echo "   💰 Hetzner VPS (€4/mo): deploy/DEPLOYMENT.md"
	@echo ""

# =========================================================================
# Docker
# =========================================================================

docker-build: ## Build Docker image
	docker build -t $(PROJECT_NAME) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(PROJECT_NAME)

docker-compose-up: ## Start with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop docker-compose
	docker-compose down

# =========================================================================
# Help
# =========================================================================

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
