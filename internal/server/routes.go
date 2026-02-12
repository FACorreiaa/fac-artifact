package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/unrolled/secure"

	"github.com/FACorreiaa/fac-artifact/assets"
	"github.com/FACorreiaa/fac-artifact/views/pages"
)

const (
	defaultCalendlyURL  = "https://calendly.com/fernandocorreia316"
	maxProposalBodySize = 1 << 20 // 1 MiB
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
	// In development mode, use more permissive CSP for templ proxy compatibility
	csp := "default-src 'self'; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; script-src 'self' 'unsafe-inline' https://unpkg.com;"
	if os.Getenv("GO_ENV") == "development" {
		// Allow templ proxy and localhost connections in development
		csp = "default-src 'self' http://localhost:* http://127.0.0.1:*; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com http://localhost:* http://127.0.0.1:*;"
	}

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
		ContentSecurityPolicy: csp,
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
	r.Get("/book-call", s.handleBookCall)
	r.Get("/proposal", s.handleProposalGet)
	r.Post("/proposal", s.handleProposalSubmit)

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

// handleBookCall redirects users to Calendly.
func (s *Server) handleBookCall(w http.ResponseWriter, r *http.Request) {
	calendlyURL := strings.TrimSpace(os.Getenv("CALENDLY_URL"))
	if calendlyURL == "" {
		calendlyURL = defaultCalendlyURL
	}
	if !strings.HasPrefix(calendlyURL, "http://") && !strings.HasPrefix(calendlyURL, "https://") {
		calendlyURL = "https://" + calendlyURL
	}

	http.Redirect(w, r, calendlyURL, http.StatusTemporaryRedirect)
}

// handleProposalGet renders the proposal request form.
func (s *Server) handleProposalGet(w http.ResponseWriter, r *http.Request) {
	s.renderProposalPage(w, r, pages.ProposalPageData{}, http.StatusOK)
}

// handleProposalSubmit validates and processes proposal form submissions.
func (s *Server) handleProposalSubmit(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxProposalBodySize)
	if err := r.ParseForm(); err != nil {
		s.renderProposalPage(w, r, pages.ProposalPageData{
			Error: "Unable to read your request. Please try again.",
		}, http.StatusBadRequest)
		return
	}

	form := pages.ProposalFormData{
		Name:        strings.TrimSpace(r.FormValue("name")),
		Email:       strings.TrimSpace(r.FormValue("email")),
		Company:     strings.TrimSpace(r.FormValue("company")),
		ProjectType: strings.TrimSpace(r.FormValue("project_type")),
		Budget:      strings.TrimSpace(r.FormValue("budget")),
		Timeline:    strings.TrimSpace(r.FormValue("timeline")),
		Details:     strings.TrimSpace(r.FormValue("details")),
	}

	// Honeypot field for basic bot filtering.
	if strings.TrimSpace(r.FormValue("website")) != "" {
		s.renderProposalPage(w, r, pages.ProposalPageData{
			Success: true,
		}, http.StatusCreated)
		return
	}

	if validationErr := validateProposalForm(form); validationErr != "" {
		s.renderProposalPage(w, r, pages.ProposalPageData{
			Form:  form,
			Error: validationErr,
		}, http.StatusBadRequest)
		return
	}

	log.Printf(
		"proposal request received name=%q email=%q company=%q project_type=%q budget=%q timeline=%q details_len=%d",
		form.Name,
		form.Email,
		form.Company,
		form.ProjectType,
		form.Budget,
		form.Timeline,
		len(form.Details),
	)

	s.renderProposalPage(w, r, pages.ProposalPageData{
		Success: true,
	}, http.StatusCreated)
}

func validateProposalForm(form pages.ProposalFormData) string {
	if len(form.Name) < 2 {
		return "Please provide your full name."
	}
	if len(form.Name) > 100 {
		return "Name is too long. Please keep it under 100 characters."
	}

	if form.Email == "" {
		return "Please provide your email address."
	}
	if _, err := mail.ParseAddress(form.Email); err != nil {
		return "Please provide a valid email address."
	}

	if form.ProjectType == "" {
		return "Please select the type of project."
	}

	if len(form.Company) > 120 {
		return "Company name is too long. Please keep it under 120 characters."
	}

	if len(form.Budget) > 60 {
		return "Budget information is too long. Please keep it concise."
	}

	if len(form.Timeline) > 60 {
		return "Timeline information is too long. Please keep it concise."
	}

	if len(form.Details) < 20 {
		return "Please share a bit more detail about your project (at least 20 characters)."
	}
	if len(form.Details) > 4000 {
		return "Project details are too long. Please keep them under 4000 characters."
	}

	return ""
}

func (s *Server) renderProposalPage(w http.ResponseWriter, r *http.Request, data pages.ProposalPageData, statusCode int) {
	w.WriteHeader(statusCode)

	component := pages.Proposal(data)
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
