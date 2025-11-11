package server

import (
	"encoding/json"
	"errors"
	"github.com/mahdi-vajdi/go-blog/internal/store"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      store.Store
}

func NewAPIServer(listenAddr string, store store.Store) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /posts", s.handleCreatePost)
	mux.HandleFunc("GET /posts", s.handleGetPosts)
	mux.HandleFunc("GET /posts/{ID}", s.handleGetPostByID)
	mux.HandleFunc("PUT /posts/{ID}", s.handleUpdatePost)
	mux.HandleFunc("DELETE /posts/{ID}", s.handleDeletePost)

	return LoggingMiddleware(mux)
}

func (s *APIServer) Run() error {
	handler := s.routes()

	log.Println("Starting server on port", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, handler)
}

func (s *APIServer) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Bad request: invalid JSON")
		return
	}

	post, err := s.store.CreatePost(body.Title, body.Content)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to created post")
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

func (s *APIServer) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := s.store.GetPosts()
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (s *APIServer) handleGetPostByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("ID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Bad request: invalid ID", http.StatusBadRequest)
		return
	}

	post, err := s.store.GetPostByID(id)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			writeJSONError(w, http.StatusNotFound, "Post not found")
			return

		} else {
			log.Printf("Failed to get post by ID %d: %v", id, err)
			writeJSONError(w, http.StatusInternalServerError, "Failed to get post")
			return
		}
	}

	writeJSON(w, http.StatusOK, post)
}

func (s *APIServer) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("ID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Bad request: invalid ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Bad request: invalid JSON")
		return
	}

	post, err := s.store.UpdatePost(id, body.Title, body.Content)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			writeJSONError(w, http.StatusNotFound, "Post not found")
			return
		} else {
			writeJSONError(w, http.StatusInternalServerError, "Failed to post post")
			return
		}
	}

	writeJSON(w, http.StatusOK, post)
}

func (s *APIServer) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("ID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Bad request: invalid ID", http.StatusBadRequest)
		return
	}

	err = s.store.DeletePost(id)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			writeJSONError(w, http.StatusNotFound, "post not found")
		} else {
			log.Printf("Failed to update post by ID %d: %v", id, err)
			writeJSONError(w, http.StatusInternalServerError, "Failed to update post")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
