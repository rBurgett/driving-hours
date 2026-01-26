package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"driving-hours/internal/auth"
	"driving-hours/internal/config"
	"driving-hours/internal/handlers"
	"driving-hours/internal/middleware"
	"driving-hours/internal/storage"
	"driving-hours/internal/templates"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage
	store, err := storage.NewJSONStorage(cfg.DataDir)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize admin on first run
	initResult, err := storage.Initialize(store, auth.HashPassword, auth.GenerateRandomPassword)
	if err != nil {
		log.Fatalf("Failed to initialize admin: %v", err)
	}

	if initResult.AdminCreated {
		fmt.Println("\n========================================")
		fmt.Println("  FIRST RUN - Admin Account Created")
		fmt.Println("========================================")
		fmt.Printf("  Email:    %s\n", initResult.AdminEmail)
		fmt.Printf("  Password: %s\n", initResult.AdminPassword)
		fmt.Println("========================================")
		fmt.Println("  Please save these credentials!")
		fmt.Println("========================================")
	}

	// Clean expired sessions
	if err := store.CleanExpiredSessions(); err != nil {
		log.Printf("Warning: Failed to clean expired sessions: %v", err)
	}

	// Initialize template renderer
	templatesDir := filepath.Join("web", "templates")
	renderer, err := templates.NewRenderer(templatesDir)
	if err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}

	// Initialize session manager
	sessions := auth.NewSessionManager(store, cfg.IsProd)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(store, sessions, renderer)
	adminHandler := handlers.NewAdminHandler(store, sessions, renderer)
	driverHandler := handlers.NewDriverHandler(store, renderer)

	// Set up router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.CSRFProtect(cfg.CSRFKey, cfg.IsProd))

	// Static files
	staticDir := filepath.Join("web", "static")
	fs := http.FileServer(http.Dir(staticDir))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, _ := sessions.GetUserFromSession(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if user.IsAdmin() {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/driver", http.StatusSeeOther)
		}
	})

	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)

	// Driver routes
	r.Route("/driver", func(r chi.Router) {
		r.Use(auth.RequireDriver(sessions))
		r.Get("/", driverHandler.Dashboard)
		r.Post("/log", driverHandler.LogHours)
		r.Get("/profile", driverHandler.Profile)
		r.Post("/profile", driverHandler.UpdateProfile)
	})

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.RequireAdmin(sessions))
		r.Get("/", adminHandler.Dashboard)
		r.Get("/users", adminHandler.ListUsers)
		r.Get("/users/new", adminHandler.NewUserForm)
		r.Post("/users", adminHandler.CreateUser)
		r.Get("/users/{id}", adminHandler.ViewDriver)
		r.Get("/users/{id}/edit", adminHandler.EditUserForm)
		r.Post("/users/{id}", adminHandler.UpdateUser)
		r.Post("/users/{id}/delete", adminHandler.DeleteUser)
		r.Get("/users/{id}/hours", adminHandler.EditHoursForm)
		r.Post("/users/{id}/hours", adminHandler.UpdateHours)
		r.Get("/profile", adminHandler.Profile)
		r.Post("/profile", adminHandler.UpdateProfile)
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
