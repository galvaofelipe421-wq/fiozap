package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type NewsletterService struct {
	sessionService *SessionService
}

func NewNewsletterService(sessionService *SessionService) *NewsletterService {
	return &NewsletterService{sessionService: sessionService}
}

func (s *NewsletterService) List(ctx context.Context, userID, sessionID string) ([]map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	newsletters, err := client.GetSubscribedNewsletters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletters: %w", err)
	}

	var result []map[string]interface{}
	for _, n := range newsletters {
		pictureURL := ""
		if n.ThreadMeta.Picture != nil {
			pictureURL = n.ThreadMeta.Picture.URL
		}
		result = append(result, map[string]interface{}{
			"jid":              n.ID.String(),
			"name":             n.ThreadMeta.Name.Text,
			"description":      n.ThreadMeta.Description.Text,
			"subscriber_count": n.ThreadMeta.SubscriberCount,
			"verification":     string(n.ThreadMeta.VerificationState),
			"picture_url":      pictureURL,
			"muted":            n.ViewerMeta != nil && n.ViewerMeta.Mute == "on",
		})
	}

	return result, nil
}

func (s *NewsletterService) GetInfo(ctx context.Context, userID, sessionID, jid string) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	newsletterJID, err := parseNewsletterJID(jid)
	if err != nil {
		return nil, err
	}

	info, err := client.GetNewsletterInfo(ctx, newsletterJID)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter info: %w", err)
	}

	return formatNewsletterInfo(info), nil
}

func (s *NewsletterService) GetInfoWithInvite(ctx context.Context, userID, sessionID, inviteKey string) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	info, err := client.GetNewsletterInfoWithInvite(ctx, inviteKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter info: %w", err)
	}

	return formatNewsletterInfo(info), nil
}

func (s *NewsletterService) GetMessages(ctx context.Context, userID, sessionID, jid string, count int, before int) ([]map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	newsletterJID, err := parseNewsletterJID(jid)
	if err != nil {
		return nil, err
	}

	params := &whatsmeow.GetNewsletterMessagesParams{
		Count: count,
	}
	if before > 0 {
		params.Before = types.MessageServerID(before)
	}

	messages, err := client.GetNewsletterMessages(ctx, newsletterJID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter messages: %w", err)
	}

	var result []map[string]interface{}
	for _, msg := range messages {
		result = append(result, map[string]interface{}{
			"server_id":   msg.MessageServerID,
			"message_id":  msg.MessageID,
			"type":        msg.Type,
			"views_count": msg.ViewsCount,
			"reactions":   msg.ReactionCounts,
			"message":     msg.Message,
			"timestamp":   msg.Timestamp.Unix(),
		})
	}

	return result, nil
}

func (s *NewsletterService) Follow(ctx context.Context, userID, sessionID, jid string) error {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return errors.New("no session")
	}

	newsletterJID, err := parseNewsletterJID(jid)
	if err != nil {
		return err
	}

	return client.FollowNewsletter(ctx, newsletterJID)
}

func (s *NewsletterService) Unfollow(ctx context.Context, userID, sessionID, jid string) error {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return errors.New("no session")
	}

	newsletterJID, err := parseNewsletterJID(jid)
	if err != nil {
		return err
	}

	return client.UnfollowNewsletter(ctx, newsletterJID)
}

func (s *NewsletterService) ToggleMute(ctx context.Context, userID, sessionID, jid string, mute bool) error {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return errors.New("no session")
	}

	newsletterJID, err := parseNewsletterJID(jid)
	if err != nil {
		return err
	}

	return client.NewsletterToggleMute(ctx, newsletterJID, mute)
}

func (s *NewsletterService) MarkViewed(ctx context.Context, userID, sessionID, jid string, serverIDs []int) error {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return errors.New("no session")
	}

	newsletterJID, err := parseNewsletterJID(jid)
	if err != nil {
		return err
	}

	msgServerIDs := make([]types.MessageServerID, len(serverIDs))
	for i, id := range serverIDs {
		msgServerIDs[i] = types.MessageServerID(id)
	}

	return client.NewsletterMarkViewed(ctx, newsletterJID, msgServerIDs)
}

func (s *NewsletterService) SendReaction(ctx context.Context, userID, sessionID, jid string, serverID int, reaction string) error {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return errors.New("no session")
	}

	newsletterJID, err := parseNewsletterJID(jid)
	if err != nil {
		return err
	}

	return client.NewsletterSendReaction(ctx, newsletterJID, types.MessageServerID(serverID), reaction, "")
}

func (s *NewsletterService) SubscribeLiveUpdates(ctx context.Context, userID, sessionID, jid string) (time.Duration, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return 0, errors.New("no session")
	}

	newsletterJID, err := parseNewsletterJID(jid)
	if err != nil {
		return 0, err
	}

	return client.NewsletterSubscribeLiveUpdates(ctx, newsletterJID)
}

func (s *NewsletterService) Create(ctx context.Context, userID, sessionID, name, description, picture string) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	params := whatsmeow.CreateNewsletterParams{
		Name:        name,
		Description: description,
	}

	if picture != "" {
		pictureData, err := base64.StdEncoding.DecodeString(picture)
		if err != nil {
			return nil, fmt.Errorf("invalid picture base64: %w", err)
		}
		params.Picture = pictureData
	}

	info, err := client.CreateNewsletter(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create newsletter: %w", err)
	}

	return formatNewsletterInfo(info), nil
}

func parseNewsletterJID(jid string) (types.JID, error) {
	if jid == "" {
		return types.JID{}, errors.New("jid is required")
	}

	if !strings.Contains(jid, "@") {
		jid = jid + "@newsletter"
	}

	parsed, err := types.ParseJID(jid)
	if err != nil {
		return types.JID{}, fmt.Errorf("invalid JID: %w", err)
	}

	return parsed, nil
}

func formatNewsletterInfo(info *types.NewsletterMetadata) map[string]interface{} {
	if info == nil {
		return nil
	}

	pictureURL := ""
	if info.ThreadMeta.Picture != nil {
		pictureURL = info.ThreadMeta.Picture.URL
	}

	result := map[string]interface{}{
		"jid":              info.ID.String(),
		"name":             info.ThreadMeta.Name.Text,
		"description":      info.ThreadMeta.Description.Text,
		"subscriber_count": info.ThreadMeta.SubscriberCount,
		"verification":     string(info.ThreadMeta.VerificationState),
		"picture_url":      pictureURL,
		"invite_code":      info.ThreadMeta.InviteCode,
		"state":            string(info.State.Type),
	}

	if info.ViewerMeta != nil {
		result["role"] = string(info.ViewerMeta.Role)
		result["muted"] = info.ViewerMeta.Mute == "on"
	}

	return result
}
