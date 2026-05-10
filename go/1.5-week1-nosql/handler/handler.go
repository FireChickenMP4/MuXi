package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/FireChickenMP4/MuXi/go/1.5-week1-nosql/models"
)

type Handler struct {
	repo models.Repository
}

func New(repo models.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 1 && parts[0] == "posts" && r.Method == http.MethodGet:
		h.listPosts(w, r)

	case len(parts) == 1 && parts[0] == "posts" && r.Method == http.MethodPost:
		h.createPost(w, r)

	case len(parts) == 2 && parts[0] == "posts" && r.Method == http.MethodGet:
		h.getPost(w, r, parts[1])

	case len(parts) == 2 && parts[0] == "posts" && r.Method == http.MethodPut:
		h.updatePost(w, r, parts[1])

	case len(parts) == 2 && parts[0] == "posts" && r.Method == http.MethodDelete:
		h.deletePost(w, r, parts[1])

	case len(parts) == 3 && parts[0] == "posts" && parts[2] == "comments" && r.Method == http.MethodPost:
		h.addComment(w, r, parts[1])

	case len(parts) == 2 && parts[0] == "comments" && r.Method == http.MethodDelete:
		h.deleteComment(w, r, parts[1])

	default:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	}
}

func (h *Handler) listPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.repo.ListPosts()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Title == "" || req.Content == "" || req.Author == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title, content and author are required"})
		return
	}
	post, err := h.repo.CreatePost(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request, id string) {
	post, err := h.repo.GetPostWithComments(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request, id string) {
	var req models.UpdatePostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	post, err := h.repo.UpdatePost(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.repo.DeletePost(id); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (h *Handler) addComment(w http.ResponseWriter, r *http.Request, postID string) {
	var req models.CreateCommentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Content == "" || req.Author == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "content and author are required"})
		return
	}
	comment, err := h.repo.AddComment(postID, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.repo.DeleteComment(id); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
