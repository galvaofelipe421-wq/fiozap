package webhook

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"fiozap/internal/database/repository"
	"fiozap/internal/logger"
)

const (
	pollInterval   = 2 * time.Second
	sendTimeout    = 10 * time.Second
	batchSize      = 50
	eventAll       = "All"
	componentName  = "webhook"
)

type Dispatcher struct {
	webhookRepo *repository.WebhookRepository
	sessionRepo *repository.SessionRepository
	sender      *Sender
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

func NewDispatcher(webhookRepo *repository.WebhookRepository, sessionRepo *repository.SessionRepository) *Dispatcher {
	return &Dispatcher{
		webhookRepo: webhookRepo,
		sessionRepo: sessionRepo,
		sender:      NewSender(),
		stopCh:      make(chan struct{}),
	}
}

func (d *Dispatcher) Start() {
	d.wg.Add(1)
	go d.processLoop()
	logger.Component(componentName).Str("status", "running").Msg("dispatcher started")
}

func (d *Dispatcher) Stop() {
	close(d.stopCh)
	d.wg.Wait()
	logger.Component(componentName).Str("status", "stopped").Msg("dispatcher stopped")
}

func (d *Dispatcher) processLoop() {
	defer d.wg.Done()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopCh:
			return
		case <-ticker.C:
			d.processPending()
		}
	}
}

func (d *Dispatcher) processPending() {
	events, err := d.webhookRepo.GetPending(batchSize)
	if err != nil {
		logger.WithError(err).Str("component", componentName).Msg("failed to get pending")
		return
	}

	for _, event := range events {
		d.processEvent(event)
	}
}

func (d *Dispatcher) processEvent(event repository.WebhookEvent) {
	if event.SessionID == "" {
		logger.WarnComponent(componentName).Int64("id", event.ID).Msg("no session_id, skipping")
		_ = d.webhookRepo.MarkFailed(event.ID)
		return
	}

	session, err := d.sessionRepo.GetByID(event.SessionID)
	if err != nil {
		logger.WarnComponent(componentName).Int64("id", event.ID).Err(err).Msg("session not found")
		_ = d.webhookRepo.MarkFailed(event.ID)
		return
	}

	if session.Webhook == "" {
		_ = d.webhookRepo.MarkFailed(event.ID)
		return
	}

	if !d.shouldSendEvent(session.Events, event.EventType) {
		_ = d.webhookRepo.MarkSent(event.ID)
		return
	}

	d.sendWebhook(event, session.Webhook)
}

func (d *Dispatcher) sendWebhook(event repository.WebhookEvent, url string) {
	var data interface{}
	_ = json.Unmarshal(event.Payload, &data)

	payload := &WebhookPayload{
		Event:     event.EventType,
		Timestamp: event.CreatedAt.Unix(),
		Data:      data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	if err := d.sender.Send(ctx, url, payload); err != nil {
		logger.WarnComponent(componentName).Int64("id", event.ID).Err(err).Msg("send failed")
		_ = d.webhookRepo.MarkFailed(event.ID)
	} else {
		logger.DebugComponent(componentName).Int64("id", event.ID).Msg("sent")
		_ = d.webhookRepo.MarkSent(event.ID)
	}
}

func (d *Dispatcher) shouldSendEvent(subscribedEvents, eventType string) bool {
	if subscribedEvents == "" {
		return false
	}

	for _, e := range strings.Split(subscribedEvents, ",") {
		if e == eventAll || e == eventType {
			return true
		}
	}
	return false
}

func (d *Dispatcher) Enqueue(userID, eventType string, data interface{}) error {
	return d.webhookRepo.Create(userID, "", eventType, data)
}

func (d *Dispatcher) EnqueueSession(userID, sessionID, eventType string, data interface{}) error {
	return d.webhookRepo.Create(userID, sessionID, eventType, data)
}
