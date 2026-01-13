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
	logger.Component("webhook").Str("status", "running").Msg("dispatcher started")
}

func (d *Dispatcher) Stop() {
	close(d.stopCh)
	d.wg.Wait()
	logger.Component("webhook").Str("status", "stopped").Msg("dispatcher stopped")
}

func (d *Dispatcher) processLoop() {
	defer d.wg.Done()

	ticker := time.NewTicker(2 * time.Second)
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
	events, err := d.webhookRepo.GetPending(50)
	if err != nil {
		logger.WithError(err).Str("component", "webhook").Msg("failed to get pending")
		return
	}

	for _, event := range events {
		if event.SessionID == "" {
			logger.Get().Warn().Str("component", "webhook").Int64("id", event.ID).Msg("no session_id, skipping")
			_ = d.webhookRepo.MarkFailed(event.ID)
			continue
		}

		session, err := d.sessionRepo.GetByID(event.SessionID)
		if err != nil {
			logger.Get().Warn().Str("component", "webhook").Int64("id", event.ID).Err(err).Msg("session not found")
			_ = d.webhookRepo.MarkFailed(event.ID)
			continue
		}

		if session.Webhook == "" {
			_ = d.webhookRepo.MarkFailed(event.ID)
			continue
		}

		if !d.shouldSendEvent(session.Events, event.EventType) {
			_ = d.webhookRepo.MarkSent(event.ID)
			continue
		}

		var data interface{}
		_ = json.Unmarshal(event.Payload, &data)

		payload := &WebhookPayload{
			Event:     event.EventType,
			Timestamp: event.CreatedAt.Unix(),
			Data:      data,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = d.sender.Send(ctx, session.Webhook, payload)
		cancel()

		if err != nil {
			logger.Get().Warn().Str("component", "webhook").Int64("id", event.ID).Err(err).Msg("send failed")
			_ = d.webhookRepo.MarkFailed(event.ID)
		} else {
			logger.Get().Debug().Str("component", "webhook").Int64("id", event.ID).Msg("sent")
			_ = d.webhookRepo.MarkSent(event.ID)
		}
	}
}

func (d *Dispatcher) shouldSendEvent(subscribedEvents, eventType string) bool {
	if subscribedEvents == "" {
		return false
	}

	events := strings.Split(subscribedEvents, ",")
	for _, e := range events {
		if e == "All" || e == eventType {
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
