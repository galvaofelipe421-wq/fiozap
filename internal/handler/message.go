package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// SendText godoc
// @Summary Send text
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.TextMessage true "Message data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/text [post]
func (h *MessageHandler) SendText(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.TextMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Message == "" {
		model.RespondBadRequest(w, errors.New("message is required"))
		return
	}

	result, err := h.messageService.SendText(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendImage godoc
// @Summary Send image
// @Description Base64 data:image/* required
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.ImageMessage true "Image data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/image [post]
func (h *MessageHandler) SendImage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ImageMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Image == "" {
		model.RespondBadRequest(w, errors.New("image is required"))
		return
	}

	result, err := h.messageService.SendImage(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendAudio godoc
// @Summary Send audio
// @Description Base64 data:audio/* required
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.AudioMessage true "Audio data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/audio [post]
func (h *MessageHandler) SendAudio(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.AudioMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Audio == "" {
		model.RespondBadRequest(w, errors.New("audio is required"))
		return
	}

	result, err := h.messageService.SendAudio(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendVideo godoc
// @Summary Send video
// @Description Base64 data:video/* required
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.VideoMessage true "Video data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/video [post]
func (h *MessageHandler) SendVideo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.VideoMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Video == "" {
		model.RespondBadRequest(w, errors.New("video is required"))
		return
	}

	result, err := h.messageService.SendVideo(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendDocument godoc
// @Summary Send document
// @Description Base64 with filename required
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.DocumentMessage true "Document data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/document [post]
func (h *MessageHandler) SendDocument(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DocumentMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Document == "" {
		model.RespondBadRequest(w, errors.New("document is required"))
		return
	}

	if req.FileName == "" {
		model.RespondBadRequest(w, errors.New("filename is required"))
		return
	}

	result, err := h.messageService.SendDocument(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendLocation godoc
// @Summary Send location
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.LocationMessage true "Location data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/location [post]
func (h *MessageHandler) SendLocation(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.LocationMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.messageService.SendLocation(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendContact godoc
// @Summary Send contact
// @Description Requires vCard format
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.ContactMessage true "Contact data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/contact [post]
func (h *MessageHandler) SendContact(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ContactMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.messageService.SendContact(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// React godoc
// @Summary React to message
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.ReactionMessage true "Reaction data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/reaction [post]
func (h *MessageHandler) React(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ReactionMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.MessageID == "" {
		model.RespondBadRequest(w, errors.New("message_id is required"))
		return
	}

	result, err := h.messageService.React(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// Delete godoc
// @Summary Delete message
// @Description Revokes own messages
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.DeleteMessage true "Delete data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/delete [post]
func (h *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DeleteMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.MessageID == "" {
		model.RespondBadRequest(w, errors.New("message_id is required"))
		return
	}

	result, err := h.messageService.Delete(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendSticker godoc
// @Summary Send sticker
// @Description WebP format required
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.StickerMessage true "Sticker data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/sticker [post]
func (h *MessageHandler) SendSticker(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.StickerMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Sticker == "" {
		model.RespondBadRequest(w, errors.New("sticker is required"))
		return
	}

	result, err := h.messageService.SendSticker(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendPoll godoc
// @Summary Send poll
// @Description Min 2 options required
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.PollMessage true "Poll data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/poll [post]
func (h *MessageHandler) SendPoll(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.PollMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Header == "" {
		model.RespondBadRequest(w, errors.New("header is required"))
		return
	}

	if len(req.Options) < 2 {
		model.RespondBadRequest(w, errors.New("at least 2 options are required"))
		return
	}

	result, err := h.messageService.SendPoll(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendList godoc
// @Summary Send list
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.ListMessage true "List data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/list [post]
func (h *MessageHandler) SendList(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ListMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" || req.ButtonText == "" || req.Title == "" || req.Desc == "" {
		model.RespondBadRequest(w, errors.New("phone, button_text, title and desc are required"))
		return
	}

	if len(req.Sections) == 0 {
		model.RespondBadRequest(w, errors.New("sections are required"))
		return
	}

	result, err := h.messageService.SendList(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendButtons godoc
// @Summary Send buttons
// @Description Experimental, 1-3 buttons
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.ButtonsMessage true "Buttons data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/buttons [post]
func (h *MessageHandler) SendButtons(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ButtonsMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Title == "" {
		model.RespondBadRequest(w, errors.New("title is required"))
		return
	}

	if len(req.Buttons) == 0 || len(req.Buttons) > 3 {
		model.RespondBadRequest(w, errors.New("1-3 buttons are required"))
		return
	}

	result, err := h.messageService.SendButtons(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// Edit godoc
// @Summary Edit message
// @Tags Messages
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.EditMessage true "Edit data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/messages/edit [post]
func (h *MessageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.EditMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.MessageID == "" {
		model.RespondBadRequest(w, errors.New("message_id is required"))
		return
	}

	if req.Body == "" {
		model.RespondBadRequest(w, errors.New("body is required"))
		return
	}

	result, err := h.messageService.EditMessage(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// MarkRead godoc
// @Summary Mark as read
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.MarkReadMessage true "Mark read data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/chat/markread [post]
func (h *MessageHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.MarkReadMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if len(req.IDs) == 0 {
		model.RespondBadRequest(w, errors.New("ids are required"))
		return
	}

	if req.ChatPhone == "" {
		model.RespondBadRequest(w, errors.New("chat_phone is required"))
		return
	}

	result, err := h.messageService.MarkRead(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SetStatusText godoc
// @Summary Set status
// @Tags Status
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.StatusTextMessage true "Status text data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/status/text [post]
func (h *MessageHandler) SetStatusText(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.StatusTextMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Body == "" {
		model.RespondBadRequest(w, errors.New("body is required"))
		return
	}

	result, err := h.messageService.SetStatusText(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// DownloadImage godoc
// @Summary Download image
// @Description Returns base64
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.DownloadMediaMessage true "Download data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/chat/downloadimage [post]
func (h *MessageHandler) DownloadImage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DownloadMediaMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	result, err := h.messageService.DownloadImage(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// DownloadVideo godoc
// @Summary Download video
// @Description Returns base64
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.DownloadMediaMessage true "Download data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/chat/downloadvideo [post]
func (h *MessageHandler) DownloadVideo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DownloadMediaMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	result, err := h.messageService.DownloadVideo(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// DownloadAudio godoc
// @Summary Download audio
// @Description Returns base64
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.DownloadMediaMessage true "Download data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/chat/downloadaudio [post]
func (h *MessageHandler) DownloadAudio(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DownloadMediaMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	result, err := h.messageService.DownloadAudio(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// DownloadDocument godoc
// @Summary Download document
// @Description Returns base64
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.DownloadMediaMessage true "Download data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/chat/downloaddocument [post]
func (h *MessageHandler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DownloadMediaMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	result, err := h.messageService.DownloadDocument(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// DownloadSticker godoc
// @Summary Download sticker
// @Description Returns base64
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.DownloadMediaMessage true "Download data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/chat/downloadsticker [post]
func (h *MessageHandler) DownloadSticker(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DownloadMediaMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	result, err := h.messageService.DownloadSticker(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// ArchiveChat godoc
// @Summary Archive chat
// @Description archive=true/false
// @Tags Chat
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param message body model.ArchiveChatMessage true "Archive data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /sessions/{sessionId}/chat/archive [post]
func (h *MessageHandler) ArchiveChat(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ArchiveChatMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.messageService.ArchiveChat(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}
