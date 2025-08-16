package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/scopophobic/scopoBlog/internal/models"
	"github.com/scopophobic/scopoBlog/internal/services"
)

type PostsHandler struct {
	DB *sql.DB
}

type AuthHandler struct {
	AuthService *services.AuthService
}

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// Login handles admin authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify password
	if !h.AuthService.VerifyPassword(loginReq.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.AuthService.GenerateJWT()
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token:   token,
		Message: "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedPost, err := services.UpdatePost(h.DB, id, &post)
	if err != nil {
		log.Printf("Error updating post %d: %v", id, err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPost)
}

func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		http.Error(w, "invalid Request body", http.StatusBadRequest)
		return
	}

	newPost, err := services.CreatePost(h.DB, &post)
	if err != nil {
		log.Printf("Error creating post in database: %v", err)
		http.Error(w, "failed to creare post : ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)

}

// GetAllPosts returns a list of all published posts.
func (h *PostsHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := services.GetAllPublishedPosts(h.DB)
	if err != nil {
		log.Printf("Error getting all posts: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetPostBySlug returns a single post by its slug.
func (h *PostsHandler) GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug") // Get the slug from the URL.

	post, err := services.GetPostBySlug(h.DB, slug)
	if err != nil {
		log.Printf("Error getting post by slug %s: %v", slug, err)
		http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
		return
	}

	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostsHandler) GetAllDraftPost(w http.ResponseWriter, r *http.Request) {
	posts, err := services.GetAllDraftPost(h.DB)
	if err != nil {
		log.Printf("Error getting all draft posts: %v", err)
		http.Error(w, "Failed to retrieve draft posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = services.DeletePost(h.DB, id)
	if err != nil {
		log.Printf("Error deleting post %d: %v", id, err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
