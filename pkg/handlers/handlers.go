package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ssych/file_service/pkg/store"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store *store.Store
}

func NewHandler(st *store.Store) *Handler {
	return &Handler{store: st}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.store.UserFindByLogin(ctx, req.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok := verifyPassword(req.Password, user.PasswordHash); !ok {
		http.Error(w, "error", http.StatusForbidden)
		return
	}

	token, err := h.store.CreateSession(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}

func (h *Handler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currentUser, ok := r.Context().Value("current_user").(int64)
	if !ok {
		http.Error(w, "current user invalid", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	name := r.PathValue("asset_name")

	if err = h.store.CreateAsset(ctx, name, currentUser, body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("success"))
}

func (h *Handler) FindAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := r.PathValue("asset_name")
	currentUser, ok := r.Context().Value("current_user").(int64)
	if !ok {
		http.Error(w, "current user invalid", http.StatusBadRequest)
		return
	}

	body, err := h.store.AssetFindByName(ctx, name, currentUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(body)
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
