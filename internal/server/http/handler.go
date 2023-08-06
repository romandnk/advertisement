package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/service"
)

type Handler struct {
	hl      *chi.Mux
	service service.Services
	logger  logger.Logger
}

func NewHandler(service service.Services, logger logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/adverts", func(r chi.Router) {
				r.Post("/", h.CreateAdvert)
			})
		})
	})

	h.hl = r

	return h.hl
}
