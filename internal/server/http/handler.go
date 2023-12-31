package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	hl        *chi.Mux
	service   service.Services
	logger    logger.Logger
	secretKey string
}

func NewHandler(service service.Services, logger logger.Logger, secretKey string) *Handler {
	return &Handler{
		service:   service,
		logger:    logger,
		secretKey: secretKey,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(h.loggingMiddleware)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/sign-up", h.SignUp)
				r.Post("/sign-in", h.SignIn)
			})

			r.Route("/adverts", func(r chi.Router) {
				r.Use(h.authorizationMiddleware)
				r.Post("/", h.CreateAdvert)
				r.Delete("/{id}", h.DeleteAdvert)
			})

			r.Group(func(r chi.Router) {
				r.Get("/adverts/{id}", h.GetAdvertByID)
			})

			r.Route("/images", func(r chi.Router) {
				r.Get("/{id}", h.GetImageByID)
			})
		})
	})

	h.hl = r

	return h.hl
}

func (h *Handler) logError(message, action string, err string) {
	h.logger.Error(message,
		zap.String("action", action),
		zap.String("error", err),
	)
}
