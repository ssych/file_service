package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ssych/file_service/pkg/store"
)

var ErrEmptyAuthHeader = errors.New("auth header is empty")

type AuthMiddleware struct {
	store *store.Store
}

func NewAuthMiddleware(st *store.Store) *AuthMiddleware {
	return &AuthMiddleware{store: st}
}

func (m *AuthMiddleware) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		session, err := m.store.SessionFindByID(context.Background(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "current_user", session.UID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", ErrEmptyAuthHeader
	}

	bearerToken := strings.Split(auth, " ")
	if len(bearerToken) != 2 {
		return "", errors.New("auth header is invalid")
	}

	return bearerToken[1], nil
}
