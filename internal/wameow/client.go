package wameow

import (
	"context"
	"fmt"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"fiozap/internal/logger"
)

// EventCallback is called when a WhatsApp event occurs
type EventCallback func(eventType string, data interface{})

// Client wraps the whatsmeow client with additional functionality
type Client struct {
	wac           *whatsmeow.Client
	userID        string
	eventCallback EventCallback
	qrCallback    func(string)
}

// NewClient creates a new WhatsApp client
func NewClient(ctx context.Context, postgresConnStr string, userID string) (*Client, error) {
	container, err := sqlstore.New(ctx, "postgres", postgresConnStr, waLogger("database"))
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlstore: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	wac := whatsmeow.NewClient(deviceStore, waLogger("whatsapp"))
	client := &Client{wac: wac, userID: userID}
	wac.AddEventHandler(client.eventHandler)

	return client, nil
}

// SetEventCallback sets the callback for WhatsApp events
func (c *Client) SetEventCallback(cb EventCallback) {
	c.eventCallback = cb
}

// SetQRCallback sets the callback for QR code events
func (c *Client) SetQRCallback(cb func(string)) {
	c.qrCallback = cb
}

// Connect establishes connection to WhatsApp
func (c *Client) Connect(ctx context.Context) error {
	if c.wac.Store.ID == nil {
		return c.connectWithQR(ctx)
	}
	return c.connectExisting()
}

// Disconnect closes the WhatsApp connection
func (c *Client) Disconnect() {
	c.wac.Disconnect()
	logger.Get().Info().Str("event", "disconnected").Msg("")
}

// GetClient returns the underlying whatsmeow client
func (c *Client) GetClient() *whatsmeow.Client {
	return c.wac
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	return c.wac.IsConnected()
}

// IsLoggedIn returns whether the client is logged in
func (c *Client) IsLoggedIn() bool {
	return c.wac.IsLoggedIn()
}

// GetJID returns the client's JID
func (c *Client) GetJID() types.JID {
	if c.wac.Store.ID != nil {
		return *c.wac.Store.ID
	}
	return types.JID{}
}

// Internal helpers

func waLogger(module string) waLog.Logger {
	return waLog.Zerolog(logger.Sub(module))
}

func (c *Client) connectWithQR(ctx context.Context) error {
	qrChan, _ := c.wac.GetQRChannel(ctx)
	if err := c.wac.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	go func() {
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, logger.Writer())
				logger.Get().Info().Str("event", "qr_code").Msg("scan to login")
				if c.qrCallback != nil {
					c.qrCallback(evt.Code)
				}
				if c.eventCallback != nil {
					c.eventCallback("QR", map[string]string{"code": evt.Code})
				}
			} else {
				logger.Get().Info().Str("event", "login").Str("status", evt.Event).Msg("")
			}
		}
	}()
	return nil
}

func (c *Client) connectExisting() error {
	if err := c.wac.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	logger.Get().Info().Str("event", "session_resumed").Msg("")
	return nil
}

// Event handler

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

func (c *Client) handleMessage(v *events.Message) {
	logger.Get().Info().Msg("event=message\n" + logger.PrettyJSON(map[string]interface{}{
		"info":    v.Info,
		"message": v.Message,
	}))
	if c.eventCallback != nil {
		c.eventCallback("Message", map[string]interface{}{
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
}

func (c *Client) handleReceipt(v *events.Receipt) {
	if c.eventCallback != nil {
		c.eventCallback("ReadReceipt", map[string]interface{}{
			"chat":       v.Chat.String(),
			"sender":     v.Sender.String(),
			"type":       string(v.Type),
			"messageIds": v.MessageIDs,
			"timestamp":  v.Timestamp.Unix(),
		})
	}
}

func (c *Client) handlePresence(v *events.Presence) {
	if c.eventCallback != nil {
		c.eventCallback("Presence", map[string]interface{}{
			"from":        v.From.String(),
			"unavailable": v.Unavailable,
			"lastSeen":    v.LastSeen.Unix(),
		})
	}
}

func (c *Client) handleChatPresence(v *events.ChatPresence) {
	if c.eventCallback != nil {
		c.eventCallback("ChatPresence", map[string]interface{}{
			"chat":   v.Chat.String(),
			"sender": v.Sender.String(),
			"state":  string(v.State),
			"media":  string(v.Media),
		})
	}
}

func (c *Client) handleConnected() {
	jid := c.wac.Store.ID.String()
	logger.Get().Info().Str("event", "connected").Str("jid", jid).Msg("")
	if c.eventCallback != nil {
		c.eventCallback("Connected", map[string]interface{}{"jid": jid})
	}
}

func (c *Client) handleDisconnected() {
	logger.Get().Warn().Str("event", "disconnected").Msg("")
	if c.eventCallback != nil {
		c.eventCallback("Disconnected", nil)
	}
}

func (c *Client) handleLoggedOut(v *events.LoggedOut) {
	reason := v.Reason.String()
	logger.Get().Warn().Str("event", "logged_out").Str("reason", reason).Msg("")
	if c.eventCallback != nil {
		c.eventCallback("LoggedOut", map[string]interface{}{"reason": reason})
	}
}

func (c *Client) handleHistorySync(v *events.HistorySync) {
	if c.eventCallback != nil {
		c.eventCallback("HistorySync", map[string]interface{}{"data": v.Data})
	}
}

func (c *Client) handleCallOffer(v *events.CallOffer) {
	if c.eventCallback != nil {
		c.eventCallback("CallOffer", map[string]interface{}{
			"from":      v.CallCreator.String(),
			"timestamp": v.Timestamp.Unix(),
			"callId":    v.CallID,
		})
	}
}

func (c *Client) handleGroupInfo(v *events.GroupInfo) {
	if c.eventCallback != nil {
		c.eventCallback("GroupInfo", map[string]interface{}{
			"jid":    v.JID.String(),
			"notify": v.Notify,
		})
	}
}

func (c *Client) handleJoinedGroup(v *events.JoinedGroup) {
	if c.eventCallback != nil {
		c.eventCallback("JoinedGroup", map[string]interface{}{
			"jid":   v.JID.String(),
			"type":  v.Type,
			"name":  v.GroupInfo.Name,
			"topic": v.GroupInfo.Topic,
		})
	}
}

// Message helpers

func getExtendedText(evt *events.Message) string {
	if evt.Message != nil && evt.Message.ExtendedTextMessage != nil {
		return evt.Message.ExtendedTextMessage.GetText()
	}
	return ""
}

func getMessageType(evt *events.Message) string {
	if evt.Message == nil {
		return "unknown"
	}
	m := evt.Message
	switch {
	case m.Conversation != nil || m.ExtendedTextMessage != nil:
		return "text"
	case m.ImageMessage != nil:
		return "image"
	case m.VideoMessage != nil:
		return "video"
	case m.AudioMessage != nil:
		return "audio"
	case m.DocumentMessage != nil:
		return "document"
	case m.StickerMessage != nil:
		return "sticker"
	case m.ContactMessage != nil:
		return "contact"
	case m.LocationMessage != nil:
		return "location"
	case m.ReactionMessage != nil:
		return "reaction"
	default:
		return "unknown"
	}
}
