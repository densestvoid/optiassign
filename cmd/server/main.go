package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"optiassign/api"
	"optiassign/db"
	"optiassign/domain"
	"optiassign/web"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize templates
	if err := web.InitTemplates(); err != nil {
		log.Fatal("Failed to initialize templates:", err)
	}

	// Connect to database
	if err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := db.NewUserRepository()

	// Initialize services
	authService := domain.NewAuthService(userRepo)

	// Initialize handlers
	authHandler := api.NewAuthHandler(authService)

	// Setup routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if user is already logged in
		session := api.GetSessionFromContext(r)
		if session != nil && session.IsValid() {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
		
		// Show login page
		data := web.PageData{
			Title: "OptiAssign - Login",
			Content: template.HTML(`
				<div class="row justify-content-center">
					<div class="col-md-6">
						<div class="card">
							<div class="card-body text-center">
								<h1 class="card-title">Welcome to OptiAssign</h1>
								<p class="card-text">Group-based, randomized, prioritized, snaking-draft item assignment</p>
								<a class="btn btn-primary btn-lg" href="/auth/google/login">
									<i class="fab fa-google me-2"></i>
									Login with Google
								</a>
							</div>
						</div>
					</div>
				</div>
			`),
		}
		web.RenderTemplate(w, "base.html", data)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		data := web.PageData{
			Title: "OptiAssign - Login",
			Content: template.HTML(`
				<div class="row justify-content-center">
					<div class="col-md-6">
						<div class="card">
							<div class="card-body text-center">
								<h1 class="card-title">Welcome to OptiAssign</h1>
								<p class="card-text">Group-based, randomized, prioritized, snaking-draft item assignment</p>
								<a class="btn btn-primary btn-lg" href="/auth/google/login">
									<i class="fab fa-google me-2"></i>
									Login with Google
								</a>
							</div>
						</div>
					</div>
				</div>
			`),
		}
		web.RenderTemplate(w, "base.html", data)
	})

	// Auth routes
	r.Get("/auth/google/login", authHandler.Login)
	r.Get("/auth/google/callback", authHandler.Callback)
	r.Get("/logout", authHandler.Logout)

	// Protected routes
	r.Route("/dashboard", func(r chi.Router) {
		r.Use(api.AuthMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			session := api.GetSessionFromContext(r)
			data := web.PageData{
				Title:   "OptiAssign - Dashboard",
				Session: session,
				Content: template.HTML(`
					<div class="row">
						<div class="col-12">
							<h1>Dashboard</h1>
							<p>Welcome, ` + session.Name + `!</p>
							<div class="row">
								<div class="col-md-6">
									<div class="card">
										<div class="card-body">
											<h5 class="card-title">Create New Group</h5>
											<p class="card-text">Start a new assignment group</p>
											<a class="btn btn-primary" href="/groups/new">Create Group</a>
										</div>
									</div>
								</div>
								<div class="col-md-6">
									<div class="card">
										<div class="card-body">
											<h5 class="card-title">My Groups</h5>
											<p class="card-text">View and manage your groups</p>
											<a class="btn btn-secondary" href="/groups">View Groups</a>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				`),
			}
			web.RenderTemplate(w, "base.html", data)
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}