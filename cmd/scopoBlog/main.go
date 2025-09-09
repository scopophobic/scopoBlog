package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/scopophobic/scopoBlog/internal/api"
	"github.com/scopophobic/scopoBlog/internal/services"
)

func main() {
	cfg, err := services.LoadConfig("./config.yaml")

	if err != nil {
		log.Fatalf("Failed to load config : %v", err)
	}

	db, err := services.InitDB(cfg.Database.Path)

	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	defer db.Close()

	log.Println("Database initialized successfully.")

	// Initialize services
	authService := services.NewAuthService(cfg.Auth.JWTSecret, cfg.Auth.AdminPasswordHash)

	// Initialize handlers
	postsHandler := &api.PostsHandler{DB: db}
	authHandler := &api.AuthHandler{AuthService: authService}
	uploadService := services.NewUploadService("./uploads", cfg.Uploads.MaxFileSize, cfg.Uploads.AllowedTypes)
	uploadHandler := &api.UploadHandler{Uploader: uploadService}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Static frontend
	r.Handle("/ui/*", http.StripPrefix("/ui/", http.FileServer(http.Dir("./web/ui"))))
	r.Get("/ui", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/ui/", http.StatusFound) })
	r.Handle("/admin-ui/*", http.StripPrefix("/admin-ui/", http.FileServer(http.Dir("./web/admin"))))
	r.Get("/admin-ui", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/admin-ui/", http.StatusFound) })

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Blog API"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Public blog routes
	r.Get("/posts", postsHandler.GetAllPosts)
	r.Get("/posts/{slug}", postsHandler.GetPostBySlug)

	// Admin routes (protected with authentication)
	r.Route("/admin", func(r chi.Router) {
		// Public admin routes (no auth required)
		r.Post("/login", authHandler.Login)

		// Protected admin routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(api.AuthMiddleware(authService))

			r.Post("/posts", postsHandler.CreatePost)
			r.Put("/posts/{id}", postsHandler.UpdatePost)
			r.Delete("/posts/{id}", postsHandler.DeletePost)
			r.Get("/posts/drafts", postsHandler.GetAllDraftPost)
			r.Post("/upload", uploadHandler.Upload)
		})
	})

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
