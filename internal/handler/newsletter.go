package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type NewsletterHandler struct {
	newsletterService *service.NewsletterService
}

func NewNewsletterHandler(newsletterService *service.NewsletterService) *NewsletterHandler {
	return &NewsletterHandler{newsletterService: newsletterService}
}

// List godoc
// @Summary List newsletters
// @Tags Newsletter
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/list [get]
func (h *NewsletterHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	result, err := h.newsletterService.List(r.Context(), user.ID, session.ID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetInfo godoc
// @Summary Get info
// @Tags Newsletter
// @Produce json
// @Param jid query string true "Newsletter JID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/info [get]
func (h *NewsletterHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	jid := r.URL.Query().Get("jid")
	if jid == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	result, err := h.newsletterService.GetInfo(r.Context(), user.ID, session.ID, jid)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetInfoWithInvite godoc
// @Summary Get info by invite
// @Tags Newsletter
// @Produce json
// @Param key query string true "Invite key/code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/info/invite [get]
func (h *NewsletterHandler) GetInfoWithInvite(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		model.RespondBadRequest(w, errors.New("key is required"))
		return
	}

	result, err := h.newsletterService.GetInfoWithInvite(r.Context(), user.ID, session.ID, key)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetMessages godoc
// @Summary Get messages
// @Description count/before for pagination
// @Tags Newsletter
// @Produce json
// @Param jid query string true "Newsletter JID"
// @Param count query int false "Number of messages to fetch" default(50)
// @Param before query int false "Fetch messages before this server ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/messages [get]
func (h *NewsletterHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	jid := r.URL.Query().Get("jid")
	if jid == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	count := 50
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil && c > 0 {
			count = c
		}
	}

	before := 0
	if beforeStr := r.URL.Query().Get("before"); beforeStr != "" {
		if b, err := strconv.Atoi(beforeStr); err == nil && b > 0 {
			before = b
		}
	}

	result, err := h.newsletterService.GetMessages(r.Context(), user.ID, session.ID, jid, count, before)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// Follow godoc
// @Summary Follow
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param request body model.NewsletterFollowRequest true "Newsletter JID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/follow [post]
func (h *NewsletterHandler) Follow(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.NewsletterFollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.JID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.newsletterService.Follow(r.Context(), user.ID, session.ID, req.JID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// Unfollow godoc
// @Summary Unfollow
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param request body model.NewsletterFollowRequest true "Newsletter JID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/unfollow [post]
func (h *NewsletterHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.NewsletterFollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.JID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.newsletterService.Unfollow(r.Context(), user.ID, session.ID, req.JID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// Mute godoc
// @Summary Mute newsletter
// @Description mute=true/false
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param request body model.NewsletterMuteRequest true "Mute request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/mute [post]
func (h *NewsletterHandler) Mute(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.NewsletterMuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.JID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.newsletterService.ToggleMute(r.Context(), user.ID, session.ID, req.JID, req.Mute)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// MarkViewed godoc
// @Summary Mark viewed
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param request body model.NewsletterMarkViewedRequest true "Mark viewed request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/markviewed [post]
func (h *NewsletterHandler) MarkViewed(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.NewsletterMarkViewedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.JID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	if len(req.ServerIDs) == 0 {
		model.RespondBadRequest(w, errors.New("server_ids is required"))
		return
	}

	err := h.newsletterService.MarkViewed(r.Context(), user.ID, session.ID, req.JID, req.ServerIDs)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// SendReaction godoc
// @Summary Send reaction
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param request body model.NewsletterReactionRequest true "Reaction request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/reaction [post]
func (h *NewsletterHandler) SendReaction(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.NewsletterReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.JID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	if req.ServerID == 0 {
		model.RespondBadRequest(w, errors.New("server_id is required"))
		return
	}

	err := h.newsletterService.SendReaction(r.Context(), user.ID, session.ID, req.JID, req.ServerID, req.Reaction)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, nil)
}

// SubscribeLiveUpdates godoc
// @Summary Live updates
// @Description Returns subscription duration
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param request body model.NewsletterFollowRequest true "Newsletter JID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/liveupdates [post]
func (h *NewsletterHandler) SubscribeLiveUpdates(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.NewsletterFollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.JID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	duration, err := h.newsletterService.SubscribeLiveUpdates(r.Context(), user.ID, session.ID, req.JID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"duration": duration.String()})
}

// Create godoc
// @Summary Create newsletter
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param request body model.NewsletterCreateRequest true "Newsletter details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/newsletter/create [post]
func (h *NewsletterHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.NewsletterCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Name == "" {
		model.RespondBadRequest(w, errors.New("name is required"))
		return
	}

	result, err := h.newsletterService.Create(r.Context(), user.ID, session.ID, req.Name, req.Description, req.Picture)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}
