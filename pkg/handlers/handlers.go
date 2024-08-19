package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ssych/file_service/pkg/render"
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
		render.Error(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.store.UserFindByLogin(ctx, req.Login)
	if err != nil && err == store.ErrNotFound {
		render.Error(w, http.StatusForbidden, errors.New("login or password wrong"))
		return
	}

	if err != nil {
		render.Error(w, http.StatusInternalServerError, err)
		return
	}

	if ok := verifyPassword(req.Password, user.PasswordHash); !ok {
		render.Error(w, http.StatusForbidden, errors.New("login or password wrong"))
		return
	}

	token, err := h.store.CreateSession(ctx, user.ID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err)
		return
	}

	render.Success(w, &LoginResponse{Token: token})
}

func (h *Handler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currentUser, ok := r.Context().Value("current_user").(int64)
	if !ok {
		render.Error(w, http.StatusBadRequest, errors.New("current user is invalid"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		render.Error(w, http.StatusBadRequest, errors.New("can't read body"))
		return
	}

	name := r.PathValue("asset_name")

	if err = h.store.CreateAsset(ctx, name, currentUser, body); err != nil {
		render.Error(w, http.StatusInternalServerError, err)
		return
	}

	render.Success(w, &CreateAssetResponse{Status: "ok"})
}

func (h *Handler) FindAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := r.PathValue("asset_name")
	currentUser, ok := r.Context().Value("current_user").(int64)
	if !ok {
		render.Error(w, http.StatusBadRequest, errors.New("current user is invalid"))
		return
	}

	body, err := h.store.AssetFindByName(ctx, name, currentUser)
	if err != nil && err == store.ErrNotFound {
		render.Error(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		render.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(body)
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
