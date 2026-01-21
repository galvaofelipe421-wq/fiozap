package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"

	"fiozap/internal/config"
	"fiozap/internal/database/repository"
	"fiozap/internal/handler"
	"fiozap/internal/middleware"
	"fiozap/internal/service"
	"fiozap/internal/webhook"
)

const requestTimeout = 60 * time.Second

type Router struct {
	mux            chi.Router
	dispatcher     *webhook.Dispatcher
	sessionService *service.SessionService
}

func New(cfg *config.Config, db *sqlx.DB) *Router {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(requestTimeout))
	r.Use(middleware.Logging)
	r.Use(corsHandler(cfg))

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)

	authMiddleware := middleware.NewAuthMiddleware(userRepo)
	adminMiddleware := middleware.NewAdminMiddleware(cfg.AdminToken)
	sessionMiddleware := middleware.NewSessionMiddleware(sessionRepo)

	sessionService := service.NewSessionService(userRepo, sessionRepo, cfg)
	sessionService.SetWebhookRepo(webhookRepo)
	dispatcher := webhook.NewDispatcher(webhookRepo, sessionRepo)
	sessionService.SetDispatcher(dispatcher)

	messageService := service.NewMessageService(sessionService)
	userService := service.NewUserService(sessionService)
	groupService := service.NewGroupService(sessionService)
	newsletterService := service.NewNewsletterService(sessionService)

	healthHandler := handler.NewHealthHandler()
	adminHandler := handler.NewAdminHandler(userRepo)
	sessionHandler := handler.NewSessionHandler(sessionService)
	messageHandler := handler.NewMessageHandler(messageService)
	userHandler := handler.NewUserHandler(userService)
	groupHandler := handler.NewGroupHandler(groupService)
	newsletterHandler := handler.NewNewsletterHandler(newsletterService)
	webhookHandler := handler.NewWebhookHandler(sessionRepo)

	// Public routes
	r.Get("/health", healthHandler.GetHealth)
	r.Mount("/swagger", httpSwagger.WrapHandler)

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(adminMiddleware.Authenticate)
		r.Get("/users", adminHandler.ListUsers)
		r.Get("/users/{id}", adminHandler.ListUsers)
		r.Post("/users", adminHandler.AddUser)
		r.Put("/users/{id}", adminHandler.EditUser)
		r.Delete("/users/{id}", adminHandler.DeleteUser)
		r.Get("/sessions", sessionHandler.AdminListAllSessions)
	})

	// API routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Get("/sessions", sessionHandler.ListSessions)
		r.Post("/sessions", sessionHandler.CreateSession)

		r.Route("/sessions/{sessionId}", func(r chi.Router) {
			r.Use(sessionMiddleware.ValidateSession)
			r.Get("/", sessionHandler.GetSession)
			r.Put("/", sessionHandler.UpdateSession)
			r.Delete("/", sessionHandler.DeleteSession)
			r.Post("/connect", sessionHandler.Connect)
			r.Post("/disconnect", sessionHandler.Disconnect)
			r.Post("/logout", sessionHandler.Logout)
			r.Get("/status", sessionHandler.GetStatus)
			r.Get("/qr", sessionHandler.GetQR)
			r.Post("/pairphone", sessionHandler.PairPhone)

			r.Route("/messages", func(r chi.Router) {
				r.Post("/text", messageHandler.SendText)
				r.Post("/image", messageHandler.SendImage)
				r.Post("/audio", messageHandler.SendAudio)
				r.Post("/video", messageHandler.SendVideo)
				r.Post("/document", messageHandler.SendDocument)
				r.Post("/location", messageHandler.SendLocation)
				r.Post("/contact", messageHandler.SendContact)
				r.Post("/reaction", messageHandler.React)
				r.Post("/delete", messageHandler.Delete)
				r.Post("/sticker", messageHandler.SendSticker)
				r.Post("/poll", messageHandler.SendPoll)
				r.Post("/list", messageHandler.SendList)
				r.Post("/buttons", messageHandler.SendButtons)
				r.Post("/edit", messageHandler.Edit)
			})

			r.Route("/chat", func(r chi.Router) {
				r.Post("/presence", userHandler.ChatPresence)
				r.Post("/markread", messageHandler.MarkRead)
				r.Post("/archive", messageHandler.ArchiveChat)
				r.Post("/downloadimage", messageHandler.DownloadImage)
				r.Post("/downloadvideo", messageHandler.DownloadVideo)
				r.Post("/downloadaudio", messageHandler.DownloadAudio)
				r.Post("/downloaddocument", messageHandler.DownloadDocument)
				r.Post("/downloadsticker", messageHandler.DownloadSticker)
			})

			r.Route("/status", func(r chi.Router) {
				r.Post("/text", messageHandler.SetStatusText)
			})

			r.Route("/call", func(r chi.Router) {
				r.Post("/reject", userHandler.RejectCall)
			})

			r.Route("/user", func(r chi.Router) {
				r.Post("/info", userHandler.GetInfo)
				r.Post("/check", userHandler.CheckUser)
				r.Post("/avatar", userHandler.GetAvatar)
				r.Get("/contacts", userHandler.GetContacts)
				r.Post("/presence", userHandler.SendPresence)
				r.Get("/newsletters", userHandler.GetNewsletters)
				r.Get("/getlid", userHandler.GetUserLID)
			})

			r.Route("/group", func(r chi.Router) {
				r.Post("/create", groupHandler.Create)
				r.Get("/list", groupHandler.List)
				r.Get("/info", groupHandler.GetInfo)
				r.Get("/invitelink", groupHandler.GetInviteLink)
				r.Post("/leave", groupHandler.Leave)
				r.Post("/updateparticipants", groupHandler.UpdateParticipants)
				r.Post("/name", groupHandler.SetName)
				r.Post("/topic", groupHandler.SetTopic)
				r.Post("/photo", groupHandler.SetPhoto)
				r.Post("/photo/remove", groupHandler.RemovePhoto)
				r.Post("/announce", groupHandler.SetAnnounce)
				r.Post("/locked", groupHandler.SetLocked)
				r.Post("/ephemeral", groupHandler.SetEphemeral)
				r.Post("/join", groupHandler.Join)
				r.Post("/inviteinfo", groupHandler.GetInviteInfo)
			})

			r.Route("/webhook", func(r chi.Router) {
				r.Get("/", webhookHandler.Get)
				r.Post("/", webhookHandler.Set)
				r.Put("/", webhookHandler.Update)
				r.Delete("/", webhookHandler.Delete)
			})

			r.Route("/newsletter", func(r chi.Router) {
				r.Get("/list", newsletterHandler.List)
				r.Get("/info", newsletterHandler.GetInfo)
				r.Get("/info/invite", newsletterHandler.GetInfoWithInvite)
				r.Get("/messages", newsletterHandler.GetMessages)
				r.Post("/follow", newsletterHandler.Follow)
				r.Post("/unfollow", newsletterHandler.Unfollow)
				r.Post("/mute", newsletterHandler.Mute)
				r.Post("/markviewed", newsletterHandler.MarkViewed)
				r.Post("/reaction", newsletterHandler.SendReaction)
				r.Post("/liveupdates", newsletterHandler.SubscribeLiveUpdates)
				r.Post("/create", newsletterHandler.Create)
			})
		})
	})

	return &Router{
		mux:            r,
		dispatcher:     dispatcher,
		sessionService: sessionService,
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.mux.ServeHTTP(w, r)
}

func (rt *Router) StartDispatcher() {
	rt.dispatcher.Start()
}

func (rt *Router) StopDispatcher() {
	rt.dispatcher.Stop()
}

func (rt *Router) GetSessionService() *service.SessionService {
	return rt.sessionService
}

func corsHandler(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := cfg.CORSOrigin
			if origin == "" {
				origin = "*"
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Token, X-Request-ID")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
