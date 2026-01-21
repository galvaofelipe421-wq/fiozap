package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

// ListSessions godoc
// @Summary List sessions
// @Tags Sessions
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions [get]
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	sessions, err := h.sessionService.GetSessionsByUser(user.ID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, sessions)
}

// CreateSession godoc
// @Summary Create session
// @Tags Sessions
// @Accept json
// @Produce json
// @Param request body model.SessionCreateRequest true "Session data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions [post]
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.SessionCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Name == "" {
		model.RespondBadRequest(w, errors.New("name is required"))
		return
	}

	session, err := h.sessionService.CreateSession(user.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondCreated(w, session)
}

// GetSession godoc
// @Summary Get session
// @Tags Sessions
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId} [get]
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		model.RespondNotFound(w, errors.New("session not found"))
		return
	}

	model.RespondOK(w, session)
}

// UpdateSession godoc
// @Summary Update session
// @Tags Sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body model.SessionUpdateRequest true "Session data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId} [put]
func (h *SessionHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		model.RespondNotFound(w, errors.New("session not found"))
		return
	}

	var req model.SessionUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	updated, err := h.sessionService.UpdateSession(session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, updated)
}

// DeleteSession godoc
// @Summary Delete session
// @Tags Sessions
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId} [delete]
func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondNotFound(w, errors.New("session not found"))
		return
	}

	if err := h.sessionService.DeleteSession(user.ID, session.ID); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// Connect godoc
// @Summary Connect session
// @Description immediate=false waits 10s for verification
// @Tags Sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body model.SessionConnectRequest false "Connection options"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/connect [post]
func (h *SessionHandler) Connect(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("session not found"))
		return
	}

	var req model.SessionConnectRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			model.RespondBadRequest(w, errors.New("could not decode payload"))
			return
		}
	}

	immediate := true
	if req.Immediate != nil {
		immediate = *req.Immediate
	}

	result, err := h.sessionService.Connect(r.Context(), user.ID, session, immediate)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// Disconnect godoc
// @Summary Disconnect session
// @Tags Sessions
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/disconnect [post]
func (h *SessionHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("session not found"))
		return
	}

	if err := h.sessionService.Disconnect(user.ID, session); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// Logout godoc
// @Summary Logout session
// @Description Clears session credentials
// @Tags Sessions
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/logout [post]
func (h *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("session not found"))
		return
	}

	if err := h.sessionService.Logout(r.Context(), user.ID, session); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// GetStatus godoc
// @Summary Get status
// @Tags Sessions
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/status [get]
func (h *SessionHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("session not found"))
		return
	}

	status := h.sessionService.GetStatus(user.ID, session)
	model.RespondOK(w, status)
}

// GetQR godoc
// @Summary Get QR code
// @Tags Sessions
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/qr [get]
func (h *SessionHandler) GetQR(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("session not found"))
		return
	}

	qr, err := h.sessionService.GetQR(user.ID, session)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"qrCode": qr})
}

// PairPhone godoc
// @Summary Pair phone
// @Description Returns 8-digit linking code
// @Tags Sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body model.PairPhoneRequest true "Phone number"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/pairphone [post]
func (h *SessionHandler) PairPhone(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("session not found"))
		return
	}

	var req model.PairPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	code, err := h.sessionService.PairPhone(r.Context(), user.ID, session, req.Phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"linkingCode": code})
}

// AdminListAllSessions godoc
// @Summary List all sessions
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security AdminKeyAuth
// @Router /admin/sessions [get]
func (h *SessionHandler) AdminListAllSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionService.GetAllSessions()
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, sessions)
}
