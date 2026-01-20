package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/unrolled/secure"

	"github.com/FACorreiaa/fac-artifact/assets"
	"github.com/FACorreiaa/fac-artifact/views/pages"
)

// RegisterRoutes sets up all routes and middleware
func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// ──────────────────────────────────────────────────────────────────
	// Core Middleware
	// ──────────────────────────────────────────────────────────────────
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Context Timeout: cancels context if request takes > 60s
	r.Use(middleware.Timeout(60 * time.Second))

	// Compression (Gzip/Deflate)
	r.Use(middleware.Compress(5))

	// ──────────────────────────────────────────────────────────────────
	// Security Middleware
	// ──────────────────────────────────────────────────────────────────

	// Secure Headers (HSTS, SSL Redirect, CSP, etc)
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{}, // Empty in dev to allow any localhost port
		AllowedHostsAreRegex:  false,
		HostsProxyHeaders:     []string{"X-Forwarded-Host"},
		SSLRedirect:           false, // Set to true in production with HTTPS
		SSLHost:               "",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; script-src 'self' 'unsafe-inline' https://unpkg.com;",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	})
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secureMiddleware.Handler(next).ServeHTTP(w, r)
		})
	})

	// CORS (Cross Origin Resource Sharing)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate Limiting (100 requests / 1 minute per IP)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// ──────────────────────────────────────────────────────────────────
	// Static Assets
	// ──────────────────────────────────────────────────────────────────
	if os.Getenv("GO_ENV") == "development" {
		// DEV: Serve from disk for hot reload
		fs := http.FileServer(http.Dir("./assets"))
		r.Handle("/assets/*", http.StripPrefix("/assets", fs))
	} else {
		// PROD: Serve from embedded binary
		fs := http.FileServer(http.FS(assets.Files))
		r.Handle("/assets/*", http.StripPrefix("/assets", fs))
	}

	// ──────────────────────────────────────────────────────────────────
	// Application Routes
	// ──────────────────────────────────────────────────────────────────

	// Health check
	r.Get("/health", s.handleHealth)

	// Pages
	r.Get("/", s.handleProjects)             // Main page is Projects
	r.Get("/projects", s.handleProjects)     // Projects page
	r.Get("/about", s.handleAbout)           // About Me page
	r.Get("/curriculum", s.handleCurriculum) // Curriculum/CV page
	r.Get("/stack", s.handleStack)           // Tech Stack page
	r.Get("/blog", s.handleBlog)             // Blog page

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/hello", s.handleAPIHello)
	})

	return r
}

// handleHealth returns service health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleProjects renders the projects page (main landing page)
func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	component := pages.Projects()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleAbout renders the about me page
func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	component := pages.About()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleCurriculum renders the CV/curriculum page
func (s *Server) handleCurriculum(w http.ResponseWriter, r *http.Request) {
	component := pages.Curriculum()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleStack renders the tech stack page
func (s *Server) handleStack(w http.ResponseWriter, r *http.Request) {
	component := pages.Stack()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleBlog renders the blog page
func (s *Server) handleBlog(w http.ResponseWriter, r *http.Request) {
	component := pages.Blog()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleAPIHello is a sample JSON API endpoint
func (s *Server) handleAPIHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello from GoForge!",
	})
}
