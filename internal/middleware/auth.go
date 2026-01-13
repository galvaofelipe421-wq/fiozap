package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"fiozap/internal/database/repository"
	"fiozap/internal/logger"
	"fiozap/internal/model"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
	bearerPrefix              = "Bearer "
	headerToken               = "Token"
	headerAuth                = "Authorization"
	queryToken                = "token"
	cacheExpiration           = 5 * time.Minute
	cacheCleanup              = 10 * time.Minute
)

var (
	errMissingToken = errors.New("missing token")
	errInvalidToken = errors.New("invalid token")
	userCache       = cache.New(cacheExpiration, cacheCleanup)
)

type AuthMiddleware struct {
	userRepo *repository.UserRepository
}

func NewAuthMiddleware(userRepo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{userRepo: userRepo}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			model.RespondUnauthorized(w, errMissingToken)
			return
		}

		user, err := m.getUser(token)
		if err != nil {
			logger.Warnf("Invalid token: %s", token)
			model.RespondUnauthorized(w, errInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) getUser(token string) (*model.User, error) {
	if cached, found := userCache.Get(token); found {
		return cached.(*model.User), nil
	}

	user, err := m.userRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	userCache.Set(token, user, cache.DefaultExpiration)
	return user, nil
}

func extractToken(r *http.Request) string {
	token := r.Header.Get(headerToken)
	if token == "" {
		token = r.Header.Get(headerAuth)
	}
	if token == "" {
		token = r.URL.Query().Get(queryToken)
	}
	return strings.TrimPrefix(token, bearerPrefix)
}

func GetUserFromContext(ctx context.Context) *model.User {
	if user, ok := ctx.Value(UserContextKey).(*model.User); ok {
		return user
	}
	return nil
}

func InvalidateUserCache(token string) { userCache.Delete(token) }

func UpdateUserCache(user *model.User) {
	userCache.Set(user.Token, user, cache.DefaultExpiration)
}
