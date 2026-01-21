package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vincent-petithory/dataurl"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"fiozap/internal/logger"
	"fiozap/internal/model"
)

type MessageService struct {
	sessionService *SessionService
}

func NewMessageService(sessionService *SessionService) *MessageService {
	return &MessageService{sessionService: sessionService}
}

func (s *MessageService) SendText(ctx context.Context, userID, sessionID string, req *model.TextMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	msg := &waE2E.Message{
		Conversation: proto.String(req.Message),
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	logger.Infof("Message sent: %s", msgID)

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendImage(ctx context.Context, userID, sessionID string, req *model.ImageMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Image, "data:image") {
		dataURL, err := dataurl.DecodeString(req.Image)
		if err != nil {
			return nil, errors.New("invalid base64 image data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("image must be base64 encoded (data:image/...)")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(filedata)
	}

	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send image: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendAudio(ctx context.Context, userID, sessionID string, req *model.AudioMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Audio, "data:audio") {
		dataURL, err := dataurl.DecodeString(req.Audio)
		if err != nil {
			return nil, errors.New("invalid base64 audio data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("audio must be base64 encoded (data:audio/...)")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaAudio)
	if err != nil {
		return nil, fmt.Errorf("failed to upload audio: %w", err)
	}

	ptt := true
	if req.PTT != nil {
		ptt = *req.PTT
	}

	mimeType := req.MimeType
	if mimeType == "" {
		if ptt {
			mimeType = "audio/ogg; codecs=opus"
		} else {
			mimeType = "audio/mpeg"
		}
	}

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
			PTT:           proto.Bool(ptt),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send audio: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendVideo(ctx context.Context, userID, sessionID string, req *model.VideoMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Video, "data:video") {
		dataURL, err := dataurl.DecodeString(req.Video)
		if err != nil {
			return nil, errors.New("invalid base64 video data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("video must be base64 encoded (data:video/...)")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaVideo)
	if err != nil {
		return nil, fmt.Errorf("failed to upload video: %w", err)
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(filedata)
	}

	msg := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send video: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendDocument(ctx context.Context, userID, sessionID string, req *model.DocumentMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Document, "data:") {
		dataURL, err := dataurl.DecodeString(req.Document)
		if err != nil {
			return nil, errors.New("invalid base64 document data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("document must be base64 encoded")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaDocument)
	if err != nil {
		return nil, fmt.Errorf("failed to upload document: %w", err)
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(filedata)
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
			FileName:      proto.String(req.FileName),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send document: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendLocation(ctx context.Context, userID, sessionID string, req *model.LocationMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(req.Latitude),
			DegreesLongitude: proto.Float64(req.Longitude),
			Name:             proto.String(req.Name),
			Address:          proto.String(req.Address),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send location: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendContact(ctx context.Context, userID, sessionID string, req *model.ContactMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	msg := &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: proto.String(req.ContactName),
			Vcard:       proto.String(req.ContactVCard),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send contact: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) React(ctx context.Context, userID, sessionID string, req *model.ReactionMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msg := &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(recipient.String()),
				ID:        proto.String(req.MessageID),
				FromMe:    proto.Bool(false),
			},
			Text: proto.String(req.Emoji),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send reaction: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
	}, nil
}

func (s *MessageService) Delete(ctx context.Context, userID, sessionID string, req *model.DeleteMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	revokeMsg := client.BuildRevoke(recipient, types.EmptyJID, req.MessageID)
	resp, err := client.SendMessage(ctx, recipient, revokeMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to delete message: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
	}, nil
}

func (s *MessageService) SendSticker(ctx context.Context, userID, sessionID string, req *model.StickerMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Sticker, "data:image") {
		dataURL, err := dataurl.DecodeString(req.Sticker)
		if err != nil {
			return nil, errors.New("invalid base64 sticker data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("sticker must be base64 encoded (data:image/webp;base64,...)")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("failed to upload sticker: %w", err)
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "image/webp"
	}

	msg := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
			PngThumbnail:  req.PngThumbnail,
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send sticker: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendPoll(ctx context.Context, userID, sessionID string, req *model.PollMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	pollMessage := client.BuildPollCreation(req.Header, req.Options, 1)
	resp, err := client.SendMessage(ctx, recipient, pollMessage, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send poll: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendList(ctx context.Context, userID, sessionID string, req *model.ListMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var sections []*waE2E.ListMessage_Section
	for _, sec := range req.Sections {
		var rows []*waE2E.ListMessage_Row
		for _, item := range sec.Rows {
			rowID := item.RowID
			if rowID == "" {
				rowID = item.Title
			}
			rows = append(rows, &waE2E.ListMessage_Row{
				RowID:       proto.String(rowID),
				Title:       proto.String(item.Title),
				Description: proto.String(item.Desc),
			})
		}
		sections = append(sections, &waE2E.ListMessage_Section{
			Title: proto.String(sec.Title),
			Rows:  rows,
		})
	}

	listMsg := &waE2E.ListMessage{
		Title:       proto.String(req.Title),
		Description: proto.String(req.Desc),
		ButtonText:  proto.String(req.ButtonText),
		ListType:    waE2E.ListMessage_SINGLE_SELECT.Enum(),
		Sections:    sections,
	}

	if req.FooterText != "" {
		listMsg.FooterText = proto.String(req.FooterText)
	}

	msg := &waE2E.Message{
		ViewOnceMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				ListMessage: listMsg,
			},
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send list: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendButtons(ctx context.Context, userID, sessionID string, req *model.ButtonsMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var buttons []*waE2E.ButtonsMessage_Button
	for _, item := range req.Buttons {
		buttons = append(buttons, &waE2E.ButtonsMessage_Button{
			ButtonID:       proto.String(item.ButtonID),
			ButtonText:     &waE2E.ButtonsMessage_Button_ButtonText{DisplayText: proto.String(item.ButtonText)},
			Type:           waE2E.ButtonsMessage_Button_RESPONSE.Enum(),
			NativeFlowInfo: &waE2E.ButtonsMessage_Button_NativeFlowInfo{},
		})
	}

	buttonsMsg := &waE2E.ButtonsMessage{
		ContentText: proto.String(req.Title),
		HeaderType:  waE2E.ButtonsMessage_EMPTY.Enum(),
		Buttons:     buttons,
	}

	msg := &waE2E.Message{
		ViewOnceMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				ButtonsMessage: buttonsMsg,
			},
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send buttons: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) EditMessage(ctx context.Context, userID, sessionID string, req *model.EditMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(req.Body),
		},
	}

	editMsg := client.BuildEdit(recipient, req.MessageID, msg)
	resp, err := client.SendMessage(ctx, recipient, editMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to edit message: %w", err)
	}

	return map[string]interface{}{
		"timestamp": resp.Timestamp.Unix(),
		"id":        req.MessageID,
	}, nil
}

func (s *MessageService) MarkRead(ctx context.Context, userID, sessionID string, req *model.MarkReadMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	chatJID, err := parseJID(req.ChatPhone)
	if err != nil {
		return nil, err
	}

	var senderJID types.JID
	if req.SenderPhone != "" {
		senderJID, err = parseJID(req.SenderPhone)
		if err != nil {
			return nil, err
		}
	}

	err = client.MarkRead(ctx, req.IDs, time.Now(), chatJID, senderJID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark as read: %w", err)
	}

	return map[string]interface{}{
	}, nil
}

func (s *MessageService) SetStatusText(ctx context.Context, userID, sessionID string, req *model.StatusTextMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	err := client.SetStatusMessage(ctx, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to set status: %w", err)
	}

	return map[string]interface{}{
	}, nil
}

func (s *MessageService) DownloadImage(ctx context.Context, userID, sessionID string, req *model.DownloadMediaMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	msg := &waE2E.Message{ImageMessage: &waE2E.ImageMessage{
		URL:           proto.String(req.URL),
		DirectPath:    proto.String(req.DirectPath),
		MediaKey:      req.MediaKey,
		Mimetype:      proto.String(req.MimeType),
		FileEncSHA256: req.FileEncSHA256,
		FileSHA256:    req.FileSHA256,
		FileLength:    &req.FileLength,
	}}

	data, err := client.Download(ctx, msg.GetImageMessage())
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	dataURL := dataurl.New(data, req.MimeType)
	return map[string]interface{}{
		"mimetype": req.MimeType,
		"data":     dataURL.String(),
	}, nil
}

func (s *MessageService) DownloadVideo(ctx context.Context, userID, sessionID string, req *model.DownloadMediaMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	msg := &waE2E.Message{VideoMessage: &waE2E.VideoMessage{
		URL:           proto.String(req.URL),
		DirectPath:    proto.String(req.DirectPath),
		MediaKey:      req.MediaKey,
		Mimetype:      proto.String(req.MimeType),
		FileEncSHA256: req.FileEncSHA256,
		FileSHA256:    req.FileSHA256,
		FileLength:    &req.FileLength,
	}}

	data, err := client.Download(ctx, msg.GetVideoMessage())
	if err != nil {
		return nil, fmt.Errorf("failed to download video: %w", err)
	}

	dataURL := dataurl.New(data, req.MimeType)
	return map[string]interface{}{
		"mimetype": req.MimeType,
		"data":     dataURL.String(),
	}, nil
}

func (s *MessageService) DownloadAudio(ctx context.Context, userID, sessionID string, req *model.DownloadMediaMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	msg := &waE2E.Message{AudioMessage: &waE2E.AudioMessage{
		URL:           proto.String(req.URL),
		DirectPath:    proto.String(req.DirectPath),
		MediaKey:      req.MediaKey,
		Mimetype:      proto.String(req.MimeType),
		FileEncSHA256: req.FileEncSHA256,
		FileSHA256:    req.FileSHA256,
		FileLength:    &req.FileLength,
	}}

	data, err := client.Download(ctx, msg.GetAudioMessage())
	if err != nil {
		return nil, fmt.Errorf("failed to download audio: %w", err)
	}

	dataURL := dataurl.New(data, req.MimeType)
	return map[string]interface{}{
		"mimetype": req.MimeType,
		"data":     dataURL.String(),
	}, nil
}

func (s *MessageService) DownloadDocument(ctx context.Context, userID, sessionID string, req *model.DownloadMediaMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	msg := &waE2E.Message{DocumentMessage: &waE2E.DocumentMessage{
		URL:           proto.String(req.URL),
		DirectPath:    proto.String(req.DirectPath),
		MediaKey:      req.MediaKey,
		Mimetype:      proto.String(req.MimeType),
		FileEncSHA256: req.FileEncSHA256,
		FileSHA256:    req.FileSHA256,
		FileLength:    &req.FileLength,
	}}

	data, err := client.Download(ctx, msg.GetDocumentMessage())
	if err != nil {
		return nil, fmt.Errorf("failed to download document: %w", err)
	}

	dataURL := dataurl.New(data, req.MimeType)
	return map[string]interface{}{
		"mimetype": req.MimeType,
		"data":     dataURL.String(),
	}, nil
}

func (s *MessageService) DownloadSticker(ctx context.Context, userID, sessionID string, req *model.DownloadMediaMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	msg := &waE2E.Message{StickerMessage: &waE2E.StickerMessage{
		URL:           proto.String(req.URL),
		DirectPath:    proto.String(req.DirectPath),
		MediaKey:      req.MediaKey,
		Mimetype:      proto.String(req.MimeType),
		FileEncSHA256: req.FileEncSHA256,
		FileSHA256:    req.FileSHA256,
		FileLength:    &req.FileLength,
	}}

	data, err := client.Download(ctx, msg.GetStickerMessage())
	if err != nil {
		return nil, fmt.Errorf("failed to download sticker: %w", err)
	}

	dataURL := dataurl.New(data, req.MimeType)
	return map[string]interface{}{
		"mimetype": req.MimeType,
		"data":     dataURL.String(),
	}, nil
}

func (s *MessageService) ArchiveChat(ctx context.Context, userID, sessionID string, req *model.ArchiveChatMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	chatJID, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	patch := appstate.BuildArchive(chatJID, req.Archive, time.Time{}, nil)
	err = client.SendAppState(ctx, patch)
	if err != nil {
		return nil, fmt.Errorf("failed to archive chat: %w", err)
	}

	return map[string]interface{}{
		"archived": req.Archive,
	}, nil
}

func parseJID(phone string) (types.JID, error) {
	if phone == "" {
		return types.JID{}, errors.New("phone is required")
	}

	if phone[0] == '+' {
		phone = phone[1:]
	}

	if !strings.Contains(phone, "@") {
		return types.NewJID(phone, types.DefaultUserServer), nil
	}

	jid, err := types.ParseJID(phone)
	if err != nil {
		return types.JID{}, fmt.Errorf("invalid JID: %w", err)
	}

	return jid, nil
}
