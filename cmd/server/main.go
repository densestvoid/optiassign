package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"optiassign/api"
	"optiassign/config"
	"optiassign/db"
	"optiassign/domain"
	"optiassign/web"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize templates
	if err := web.InitTemplates(); err != nil {
		log.Fatal("Failed to initialize templates:", err)
	}

	// Connect to database
	if err := db.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations using Goose directly
	// Note: In production, migrations should be run separately
	// For development, we'll skip auto-migration
	log.Println("Note: Run 'task db-migrate' to apply database migrations")

	// Initialize repositories
	userRepo := db.NewUserRepository()
	groupRepo := db.NewGroupRepository()
	itemRepo := db.NewItemRepository()
	participantRepo := db.NewParticipantRepository()
	prioritizationRepo := db.NewPrioritizationRepository()
	assignmentRepo := db.NewAssignmentRepository()

	// Initialize services
	authService := domain.NewAuthService(userRepo)
	groupService := domain.NewGroupService(groupRepo, itemRepo, participantRepo, userRepo)
	assignmentAlgorithm := domain.NewAssignmentAlgorithm(groupRepo, itemRepo, participantRepo, prioritizationRepo, assignmentRepo)
	emailService := domain.NewEmailService(cfg.EmailLogToConsole)

	// Initialize handlers
	authHandler := api.NewAuthHandler(authService, cfg)
	groupHandler := api.NewGroupHandler(groupService, userRepo, assignmentAlgorithm, emailService)
	participantHandler := api.NewParticipantHandler(groupService, prioritizationRepo, participantRepo, itemRepo)

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

	// Group routes
	r.Route("/groups", func(r chi.Router) {
		r.Use(api.AuthMiddleware)
		r.Get("/", groupHandler.ListGroups)
		r.Get("/new", groupHandler.NewGroupForm)
		r.Post("/", groupHandler.CreateGroup)
		r.Get("/{id}", groupHandler.ViewGroup)
		r.Post("/{id}/execute", groupHandler.ExecuteAssignment)
	})

	// Participant routes (no auth required - uses tokens)
	r.Route("/participant", func(r chi.Router) {
		r.Get("/{token}", participantHandler.ParticipantAccess)
		r.Get("/{token}/priorities", participantHandler.PriorityForm)
		r.Post("/{token}/priorities", participantHandler.SubmitPriorities)
	})

	// Task management is handled by Taskfile.yml

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}