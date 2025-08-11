package main

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// "internal/service"
	// "github.com/scopo-user/scopoBlog/internal/services"
	"github.com/scopophobic/scopoBlog/internal/services"
	"github.com/scopophobic/scopoBlog/internal/api"
)


func main(){
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	postsHandler := &api.PostsHandler{DB: db}


	// api routes 

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Blog API"))
	})

	r.Route("/admin", func(r chi.Router){
		r.Post("/posts",postsHandler.CreatePost)
		r.Put("/posts/{id}", postsHandler.UpdatePost)
	})
	

	// public routes
	r.Get("/posts", postsHandler.GetAllPosts)
	r.Get("/posts/{slug}", postsHandler.GetPostBySlug)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Println("Starting server on :8080")

	// 💡 Use your chi router 'r' here instead of 'nil'
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}


}