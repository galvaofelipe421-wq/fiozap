package wameow

import (
	"go.mau.fi/whatsmeow/types/events"

	"fiozap/internal/logger"
)

const (
	eventQR           = "QR"
	eventMessage      = "Message"
	eventReadReceipt  = "ReadReceipt"
	eventPresence     = "Presence"
	eventChatPresence = "ChatPresence"
	eventConnected    = "Connected"
	eventDisconnected = "Disconnected"
	eventLoggedOut    = "LoggedOut"
	eventHistorySync  = "HistorySync"
	eventCallOffer    = "CallOffer"
	eventGroupInfo    = "GroupInfo"
	eventJoinedGroup  = "JoinedGroup"
)

func (c *Client) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		c.handleMessage(v)
	case *events.Receipt:
		c.handleReceipt(v)
	case *events.Presence:
		c.handlePresence(v)
	case *events.ChatPresence:
		c.handleChatPresence(v)
	case *events.Connected:
		c.handleConnected()
	case *events.Disconnected:
		c.handleDisconnected()
	case *events.LoggedOut:
		c.handleLoggedOut(v)
	case *events.HistorySync:
		c.handleHistorySync(v)
	case *events.CallOffer:
		c.handleCallOffer(v)
	case *events.GroupInfo:
		c.handleGroupInfo(v)
	case *events.JoinedGroup:
		c.handleJoinedGroup(v)
	}
}

func (c *Client) emit(eventType string, data interface{}) {
	if c.eventCallback != nil {
		c.eventCallback(eventType, data)
	}
}

func (c *Client) handleMessage(v *events.Message) {
	logger.Get().Info().Msg("event=message\n" + logger.PrettyJSON(map[string]interface{}{
		"info":    v.Info,
		"message": v.Message,
	}))

	c.emit(eventMessage, map[string]interface{}{
		"from":         v.Info.Sender.String(),
		"chat":         v.Info.Chat.String(),
		"id":           v.Info.ID,
		"timestamp":    v.Info.Timestamp.Unix(),
		"pushName":     v.Info.PushName,
		"isGroup":      v.Info.IsGroup,
		"isFromMe":     v.Info.IsFromMe,
		"type":         getMessageType(v),
		"text":         v.Message.GetConversation(),
		"extendedText": getExtendedText(v),
	})
}

func (c *Client) handleReceipt(v *events.Receipt) {
	c.emit(eventReadReceipt, map[string]interface{}{
		"chat":       v.Chat.String(),
		"sender":     v.Sender.String(),
		"type":       string(v.Type),
		"messageIds": v.MessageIDs,
		"timestamp":  v.Timestamp.Unix(),
	})
}

func (c *Client) handlePresence(v *events.Presence) {
	c.emit(eventPresence, map[string]interface{}{
		"from":        v.From.String(),
		"unavailable": v.Unavailable,
		"lastSeen":    v.LastSeen.Unix(),
	})
}

func (c *Client) handleChatPresence(v *events.ChatPresence) {
	c.emit(eventChatPresence, map[string]interface{}{
		"chat":   v.Chat.String(),
		"sender": v.Sender.String(),
		"state":  string(v.State),
		"media":  string(v.Media),
	})
}

func (c *Client) handleConnected() {
	jid := c.wac.Store.ID.String()
	logger.Get().Info().Str("event", "connected").Str("jid", jid).Msg("")
	c.emit(eventConnected, map[string]interface{}{"jid": jid})
}

func (c *Client) handleDisconnected() {
	logger.Get().Warn().Str("event", "disconnected").Msg("")
	c.emit(eventDisconnected, nil)
}

func (c *Client) handleLoggedOut(v *events.LoggedOut) {
	reason := v.Reason.String()
	logger.Get().Warn().Str("event", "logged_out").Str("reason", reason).Msg("")
	c.emit(eventLoggedOut, map[string]interface{}{"reason": reason})
}

func (c *Client) handleHistorySync(v *events.HistorySync) {
	c.emit(eventHistorySync, map[string]interface{}{"data": v.Data})
}

func (c *Client) handleCallOffer(v *events.CallOffer) {
	c.emit(eventCallOffer, map[string]interface{}{
		"from":      v.CallCreator.String(),
		"timestamp": v.Timestamp.Unix(),
		"callId":    v.CallID,
	})
}

func (c *Client) handleGroupInfo(v *events.GroupInfo) {
	c.emit(eventGroupInfo, map[string]interface{}{
		"jid":    v.JID.String(),
		"notify": v.Notify,
	})
}

func (c *Client) handleJoinedGroup(v *events.JoinedGroup) {
	c.emit(eventJoinedGroup, map[string]interface{}{
		"jid":   v.JID.String(),
		"type":  v.Type,
		"name":  v.Name,
		"topic": v.Topic,
	})
}

func getExtendedText(evt *events.Message) string {
	if evt.Message != nil && evt.Message.ExtendedTextMessage != nil {
		return evt.Message.ExtendedTextMessage.GetText()
	}
	return ""
}

func getMessageType(evt *events.Message) string {
	if evt.Message == nil {
		return msgTypeUnknown
	}
	m := evt.Message
	switch {
	case m.Conversation != nil || m.ExtendedTextMessage != nil:
		return msgTypeText
	case m.ImageMessage != nil:
		return msgTypeImage
	case m.VideoMessage != nil:
		return msgTypeVideo
	case m.AudioMessage != nil:
		return msgTypeAudio
	case m.DocumentMessage != nil:
		return msgTypeDocument
	case m.StickerMessage != nil:
		return msgTypeSticker
	case m.ContactMessage != nil:
		return msgTypeContact
	case m.LocationMessage != nil:
		return msgTypeLocation
	case m.ReactionMessage != nil:
		return msgTypeReaction
	default:
		return msgTypeUnknown
	}
}
