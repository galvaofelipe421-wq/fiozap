package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// Create godoc
// @Summary Create group
// @Description Create a new WhatsApp group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{name=string,participants=[]string} true "Group data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/create [post]
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Name         string   `json:"name"`
		Participants []string `json:"participants"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Name == "" {
		model.RespondBadRequest(w, errors.New("name is required"))
		return
	}

	result, err := h.groupService.Create(r.Context(), user.ID, session.ID, req.Name, req.Participants)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// List godoc
// @Summary List groups
// @Description List all joined WhatsApp groups
// @Tags Group
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/list [get]
func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	result, err := h.groupService.List(r.Context(), user.ID, session.ID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetInfo godoc
// @Summary Get group info
// @Description Get detailed information about a group
// @Tags Group
// @Produce json
// @Param jid query string true "Group JID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/info [get]
func (h *GroupHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	groupJID := r.URL.Query().Get("jid")
	if groupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	result, err := h.groupService.GetInfo(r.Context(), user.ID, session.ID, groupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetInviteLink godoc
// @Summary Get group invite link
// @Description Get the invite link for a group
// @Tags Group
// @Produce json
// @Param jid query string true "Group JID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/invitelink [get]
func (h *GroupHandler) GetInviteLink(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	groupJID := r.URL.Query().Get("jid")
	if groupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	result, err := h.groupService.GetInviteLink(r.Context(), user.ID, session.ID, groupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// Leave godoc
// @Summary Leave group
// @Description Leave a WhatsApp group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string} true "Group JID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/leave [post]
func (h *GroupHandler) Leave(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID string `json:"jid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.groupService.Leave(r.Context(), user.ID, session.ID, req.GroupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Left group"})
}

// UpdateParticipants godoc
// @Summary Update group participants
// @Description Add, remove, promote or demote group participants
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string,participants=[]string,action=string} true "Participants data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/updateparticipants [post]
func (h *GroupHandler) UpdateParticipants(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID     string   `json:"jid"`
		Participants []string `json:"participants"`
		Action       string   `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	if len(req.Participants) == 0 {
		model.RespondBadRequest(w, errors.New("participants is required"))
		return
	}

	if req.Action == "" {
		req.Action = "add"
	}

	result, err := h.groupService.UpdateParticipants(r.Context(), user.ID, session.ID, req.GroupJID, req.Participants, req.Action)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SetName godoc
// @Summary Set group name
// @Description Update the name of a group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string,name=string} true "Group name data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/name [post]
func (h *GroupHandler) SetName(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID string `json:"jid"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" || req.Name == "" {
		model.RespondBadRequest(w, errors.New("jid and name are required"))
		return
	}

	err := h.groupService.SetName(r.Context(), user.ID, session.ID, req.GroupJID, req.Name)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group name updated"})
}

// SetTopic godoc
// @Summary Set group topic
// @Description Update the topic/description of a group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string,topic=string} true "Group topic data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/topic [post]
func (h *GroupHandler) SetTopic(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID string `json:"jid"`
		Topic    string `json:"topic"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.groupService.SetTopic(r.Context(), user.ID, session.ID, req.GroupJID, req.Topic)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group topic updated"})
}

// SetPhoto godoc
// @Summary Set group photo
// @Description Set the photo of a group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body model.GroupPhotoRequest true "Group photo data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/photo [post]
func (h *GroupHandler) SetPhoto(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.GroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	if req.Image == "" {
		model.RespondBadRequest(w, errors.New("image is required"))
		return
	}

	// Decode base64 image
	var imageData []byte
	if len(req.Image) > 10 && req.Image[0:10] == "data:image" {
		// Find the base64 part
		idx := 0
		for i, c := range req.Image {
			if c == ',' {
				idx = i + 1
				break
			}
		}
		if idx > 0 {
			decoded := make([]byte, len(req.Image[idx:]))
			n, _ := base64Decode(decoded, []byte(req.Image[idx:]))
			imageData = decoded[:n]
		}
	}

	if len(imageData) == 0 {
		model.RespondBadRequest(w, errors.New("invalid image data"))
		return
	}

	pictureID, err := h.groupService.SetPhoto(r.Context(), user.ID, session.ID, req.GroupJID, imageData)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group photo set", "picture_id": pictureID})
}

// RemovePhoto godoc
// @Summary Remove group photo
// @Description Remove the photo of a group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string} true "Group JID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/photo/remove [post]
func (h *GroupHandler) RemovePhoto(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID string `json:"jid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.groupService.RemovePhoto(r.Context(), user.ID, session.ID, req.GroupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group photo removed"})
}

// SetAnnounce godoc
// @Summary Set group announce mode
// @Description Set the announce mode of a group (only admins can send messages)
// @Tags Group
// @Accept json
// @Produce json
// @Param request body model.GroupAnnounceRequest true "Group announce data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/announce [post]
func (h *GroupHandler) SetAnnounce(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.GroupAnnounceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.groupService.SetAnnounce(r.Context(), user.ID, session.ID, req.GroupJID, req.Announce)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group announce mode updated"})
}

// SetLocked godoc
// @Summary Set group locked mode
// @Description Set the locked mode of a group (only admins can edit group info)
// @Tags Group
// @Accept json
// @Produce json
// @Param request body model.GroupLockedRequest true "Group locked data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/locked [post]
func (h *GroupHandler) SetLocked(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.GroupLockedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.groupService.SetLocked(r.Context(), user.ID, session.ID, req.GroupJID, req.Locked)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group locked mode updated"})
}

// SetEphemeral godoc
// @Summary Set group ephemeral timer
// @Description Set the disappearing messages timer for a group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body model.GroupEphemeralRequest true "Group ephemeral data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/ephemeral [post]
func (h *GroupHandler) SetEphemeral(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.GroupEphemeralRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	if req.Duration == "" {
		model.RespondBadRequest(w, errors.New("duration is required (24h, 7d, 90d, or off)"))
		return
	}

	err := h.groupService.SetEphemeral(r.Context(), user.ID, session.ID, req.GroupJID, req.Duration)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group ephemeral timer updated"})
}

// Join godoc
// @Summary Join group via invite link
// @Description Join a group using an invite link/code
// @Tags Group
// @Accept json
// @Produce json
// @Param request body model.GroupJoinRequest true "Group join data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/join [post]
func (h *GroupHandler) Join(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.GroupJoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Code == "" {
		model.RespondBadRequest(w, errors.New("code is required"))
		return
	}

	result, err := h.groupService.Join(r.Context(), user.ID, session.ID, req.Code)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetInviteInfo godoc
// @Summary Get group invite info
// @Description Get information about a group from an invite link
// @Tags Group
// @Accept json
// @Produce json
// @Param request body model.GroupInviteInfoRequest true "Group invite info request"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Param sessionId path string true "Session ID"
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/group/inviteinfo [post]
func (h *GroupHandler) GetInviteInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.GroupInviteInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Code == "" {
		model.RespondBadRequest(w, errors.New("code is required"))
		return
	}

	result, err := h.groupService.GetInviteInfo(r.Context(), user.ID, session.ID, req.Code)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func base64Decode(dst, src []byte) (n int, err error) {
	// Simple base64 decode using encoding/base64
	import64 := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	decodeMap := make([]byte, 256)
	for i := range decodeMap {
		decodeMap[i] = 0xFF
	}
	for i := 0; i < len(import64); i++ {
		decodeMap[import64[i]] = byte(i)
	}

	si, di := 0, 0
	for si < len(src) {
		var dbuf [4]byte
		dlen := 4
		for j := 0; j < 4; j++ {
			if si >= len(src) {
				if j == 0 {
					return di, nil
				}
				dlen = j
				break
			}
			in := src[si]
			si++
			if in == '=' {
				dlen = j
				break
			}
			dbuf[j] = decodeMap[in]
		}

		val := uint(dbuf[0])<<18 | uint(dbuf[1])<<12 | uint(dbuf[2])<<6 | uint(dbuf[3])
		switch dlen {
		case 4:
			dst[di+2] = byte(val)
			fallthrough
		case 3:
			dst[di+1] = byte(val >> 8)
			fallthrough
		case 2:
			dst[di] = byte(val >> 16)
		}
		di += dlen - 1
	}
	return di, nil
}
