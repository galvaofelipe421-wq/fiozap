package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetInfo godoc
// @Summary Get user info
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{phone=[]string} true "Phone numbers"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/user/info [post]
func (h *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Phone []string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if len(req.Phone) == 0 {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.userService.GetInfo(r.Context(), user.ID, session.ID, req.Phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// CheckUser godoc
// @Summary Check users
// @Description Verifies WhatsApp registration
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{phone=[]string} true "Phone numbers"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/user/check [post]
func (h *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Phone []string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if len(req.Phone) == 0 {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.userService.CheckUser(r.Context(), user.ID, session.ID, req.Phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetAvatar godoc
// @Summary Get avatar
// @Description preview=true for thumbnail
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{phone=string,preview=bool} true "Phone and preview option"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/user/avatar [post]
func (h *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Phone   string `json:"phone"`
		Preview bool   `json:"preview"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.userService.GetAvatar(r.Context(), user.ID, session.ID, req.Phone, req.Preview)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetContacts godoc
// @Summary Get contacts
// @Tags User
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/user/contacts [get]
func (h *UserHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	result, err := h.userService.GetContacts(r.Context(), user.ID, session.ID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendPresence godoc
// @Summary Send presence
// @Description available/unavailable
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{presence=string} true "Presence (available/unavailable)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/user/presence [post]
func (h *UserHandler) SendPresence(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Presence string `json:"presence"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Presence == "" {
		req.Presence = "available"
	}

	err := h.userService.SendPresence(r.Context(), user.ID, session.ID, req.Presence)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// ChatPresence godoc
// @Summary Chat presence
// @Description state: composing/paused
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{phone=string,state=string,media=string} true "Chat presence data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/chat/presence [post]
func (h *UserHandler) ChatPresence(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Phone string `json:"phone"`
		State string `json:"state"`
		Media string `json:"media"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	err := h.userService.ChatPresence(r.Context(), user.ID, session.ID, req.Phone, req.State, req.Media)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// RejectCall godoc
// @Summary Reject call
// @Tags User
// @Accept json
// @Produce json
// @Param request body model.RejectCallRequest true "Call rejection data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/call/reject [post]
func (h *UserHandler) RejectCall(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.RejectCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.CallFrom == "" {
		model.RespondBadRequest(w, errors.New("call_from is required"))
		return
	}

	if req.CallID == "" {
		model.RespondBadRequest(w, errors.New("call_id is required"))
		return
	}

	err := h.userService.RejectCall(r.Context(), user.ID, session.ID, req.CallFrom, req.CallID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// GetNewsletters godoc
// @Summary Get newsletters
// @Tags User
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/user/newsletters [get]
func (h *UserHandler) GetNewsletters(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	result, err := h.userService.GetNewsletters(r.Context(), user.ID, session.ID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetUserLID godoc
// @Summary Get LID
// @Tags User
// @Produce json
// @Param phone query string true "Phone number"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/user/getlid [get]
func (h *UserHandler) GetUserLID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	phone := r.URL.Query().Get("phone")
	if phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.userService.GetUserLID(r.Context(), user.ID, session.ID, phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}
