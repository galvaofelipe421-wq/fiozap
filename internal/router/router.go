package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"

	"fiozap/internal/config"
	"fiozap/internal/database/repository"
	"fiozap/internal/handler"
	"fiozap/internal/middleware"
	"fiozap/internal/service"
	"fiozap/internal/webhook"
)

type Router struct {
	mux            chi.Router
	dispatcher     *webhook.Dispatcher
	sessionService *service.SessionService
}

func New(cfg *config.Config, db *sqlx.DB) *Router {
	r := chi.NewRouter()

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)

	authMiddleware := middleware.NewAuthMiddleware(userRepo)
	adminMiddleware := middleware.NewAdminMiddleware(cfg.AdminToken)
	sessionMiddleware := middleware.NewSessionMiddleware(sessionRepo)

	healthHandler := handler.NewHealthHandler()
	adminHandler := handler.NewAdminHandler(userRepo)

	sessionService := service.NewSessionService(userRepo, sessionRepo, cfg)
	sessionService.SetWebhookRepo(webhookRepo)

	dispatcher := webhook.NewDispatcher(webhookRepo, sessionRepo)
	sessionService.SetDispatcher(dispatcher)

	sessionHandler := handler.NewSessionHandler(sessionService)

	messageService := service.NewMessageService(sessionService)
	messageHandler := handler.NewMessageHandler(messageService)

	userService := service.NewUserService(sessionService)
	userHandler := handler.NewUserHandler(userService)

	groupService := service.NewGroupService(sessionService)
	groupHandler := handler.NewGroupHandler(groupService)

	webhookHandler := handler.NewWebhookHandler(sessionRepo)

	r.Use(cors)
	r.Use(middleware.Logging)

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

	// API routes (authenticated)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		// Sessions CRUD (user manages own sessions)
		r.Get("/sessions", sessionHandler.ListSessions)
		r.Post("/sessions", sessionHandler.CreateSession)

		// Session-specific routes (require session validation)
		r.Route("/sessions/{sessionId}", func(r chi.Router) {
			r.Use(sessionMiddleware.ValidateSession)

			// Session management
			r.Get("/", sessionHandler.GetSession)
			r.Put("/", sessionHandler.UpdateSession)
			r.Delete("/", sessionHandler.DeleteSession)

			// Session connection
			r.Post("/connect", sessionHandler.Connect)
			r.Post("/disconnect", sessionHandler.Disconnect)
			r.Post("/logout", sessionHandler.Logout)
			r.Get("/status", sessionHandler.GetStatus)
			r.Get("/qr", sessionHandler.GetQR)
			r.Post("/pairphone", sessionHandler.PairPhone)

			// Messages (per session)
			r.Post("/messages/text", messageHandler.SendText)
			r.Post("/messages/image", messageHandler.SendImage)
			r.Post("/messages/audio", messageHandler.SendAudio)
			r.Post("/messages/video", messageHandler.SendVideo)
			r.Post("/messages/document", messageHandler.SendDocument)
			r.Post("/messages/location", messageHandler.SendLocation)
			r.Post("/messages/contact", messageHandler.SendContact)
			r.Post("/messages/reaction", messageHandler.React)
			r.Post("/messages/delete", messageHandler.Delete)

			// User operations (per session)
			r.Post("/user/info", userHandler.GetInfo)
			r.Post("/user/check", userHandler.CheckUser)
			r.Post("/user/avatar", userHandler.GetAvatar)
			r.Get("/user/contacts", userHandler.GetContacts)
			r.Post("/user/presence", userHandler.SendPresence)
			r.Post("/chat/presence", userHandler.ChatPresence)

			// Group operations (per session)
			r.Post("/group/create", groupHandler.Create)
			r.Get("/group/list", groupHandler.List)
			r.Get("/group/info", groupHandler.GetInfo)
			r.Get("/group/invitelink", groupHandler.GetInviteLink)
			r.Post("/group/leave", groupHandler.Leave)
			r.Post("/group/updateparticipants", groupHandler.UpdateParticipants)
			r.Post("/group/name", groupHandler.SetName)
			r.Post("/group/topic", groupHandler.SetTopic)

			// Webhook (per session)
			r.Get("/webhook", webhookHandler.Get)
			r.Post("/webhook", webhookHandler.Set)
			r.Put("/webhook", webhookHandler.Update)
			r.Delete("/webhook", webhookHandler.Delete)
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

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Token")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
