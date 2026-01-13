package wameow

import (
	"context"
	"fmt"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"

	"fiozap/internal/logger"
)

const (
	driverPostgres = "postgres"
	qrEventCode    = "code"
	qrEventTimeout = "timeout"

	msgTypeText     = "text"
	msgTypeImage    = "image"
	msgTypeVideo    = "video"
	msgTypeAudio    = "audio"
	msgTypeDocument = "document"
	msgTypeSticker  = "sticker"
	msgTypeContact  = "contact"
	msgTypeLocation = "location"
	msgTypeReaction = "reaction"
	msgTypeUnknown  = "unknown"
)

type EventCallback func(eventType string, data interface{})

type Client struct {
	wac           *whatsmeow.Client
	userID        string
	eventCallback EventCallback
	qrCallback    func(string)
}

func NewClient(ctx context.Context, postgresConnStr string, userID string) (*Client, error) {
	container, err := sqlstore.New(ctx, driverPostgres, postgresConnStr, waLogger("database"))
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

func (c *Client) SetEventCallback(cb EventCallback) { c.eventCallback = cb }
func (c *Client) SetQRCallback(cb func(string))     { c.qrCallback = cb }
func (c *Client) GetClient() *whatsmeow.Client      { return c.wac }
func (c *Client) IsConnected() bool                 { return c.wac.IsConnected() }
func (c *Client) IsLoggedIn() bool                  { return c.wac.IsLoggedIn() }

func (c *Client) GetJID() types.JID {
	if c.wac.Store.ID != nil {
		return *c.wac.Store.ID
	}
	return types.JID{}
}

func (c *Client) Connect(ctx context.Context) error {
	if c.wac.Store.ID == nil {
		return c.connectWithQR(ctx)
	}
	return c.connectExisting()
}

func (c *Client) Disconnect() {
	c.wac.Disconnect()
	logger.Get().Info().Str("event", "disconnected").Msg("")
}

func waLogger(module string) waLog.Logger {
	return waLog.Zerolog(logger.Sub(module))
}

func (c *Client) connectWithQR(ctx context.Context) error {
	qrChan, _ := c.wac.GetQRChannel(context.Background())
	if err := c.wac.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	go c.handleQRChannel(qrChan)
	return nil
}

func (c *Client) handleQRChannel(qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		if evt.Event == qrEventCode {
			logger.Get().Info().Str("event", "qr_code").Msg("scan to login")
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, logger.RawWriter())
			if c.qrCallback != nil {
				c.qrCallback(evt.Code)
			}
			c.emit(eventQR, map[string]string{"code": evt.Code})
		} else {
			logger.Get().Info().Str("event", "login").Str("status", evt.Event).Msg("")
			if evt.Event == qrEventTimeout {
				logger.Get().Warn().Msg("QR code expired, please reconnect to get a new one")
			}
		}
	}
}

func (c *Client) connectExisting() error {
	if err := c.wac.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	logger.Get().Info().Str("event", "session_resumed").Msg("")
	return nil
}
