package middleware

import (
	"errors"
	"net/http"
	"strings"

	"fiozap/internal/model"
)

var (
	errMissingAdmin = errors.New("missing admin token")
	errInvalidAdmin = errors.New("invalid admin token")
)

type AdminMiddleware struct {
	adminToken string
}

func NewAdminMiddleware(adminToken string) *AdminMiddleware {
	return &AdminMiddleware{adminToken: adminToken}
}

func (m *AdminMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get(headerAuth), bearerPrefix)

		if token == "" {
			model.RespondUnauthorized(w, errMissingAdmin)
			return
		}

		if token != m.adminToken {
			model.RespondUnauthorized(w, errInvalidAdmin)
			return
		}

		next.ServeHTTP(w, r)
	})
}
