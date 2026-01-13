package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"fiozap/internal/database/repository"
	"fiozap/internal/model"
)

const (
	SessionContextKey contextKey = "session"
	paramSessionID    string     = "sessionId"
)

var (
	errUserNotFound    = errors.New("user not found")
	errSessionRequired = errors.New("session name is required")
	errSessionNotFound = errors.New("session not found")
)

type SessionMiddleware struct {
	sessionRepo *repository.SessionRepository
}

func NewSessionMiddleware(sessionRepo *repository.SessionRepository) *SessionMiddleware {
	return &SessionMiddleware{sessionRepo: sessionRepo}
}

func (m *SessionMiddleware) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			model.RespondUnauthorized(w, errUserNotFound)
			return
		}

		sessionName := chi.URLParam(r, paramSessionID)
		if sessionName == "" {
			model.RespondBadRequest(w, errSessionRequired)
			return
		}

		session, err := m.sessionRepo.GetByUserAndName(user.ID, sessionName)
		if err != nil {
			model.RespondNotFound(w, errSessionNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSessionFromContext(ctx context.Context) *model.Session {
	if session, ok := ctx.Value(SessionContextKey).(*model.Session); ok {
		return session
	}
	return nil
}
