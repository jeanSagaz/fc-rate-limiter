package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jeanSagaz/rate-limiter/internal/application/dto"
	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis/repository"
)

type Handler struct {
	IRedisRepository   repository.IRedisRepository
	TokenConfiguration []dto.TokenConfiguration
	NumberRequests     int
	Seconds            int
}

func NewHandler(iRedisRepository repository.IRedisRepository,
	tokenConfiguration []dto.TokenConfiguration,
	numberRequests int,
	seconds int) *Handler {
	return &Handler{
		IRedisRepository:   iRedisRepository,
		TokenConfiguration: tokenConfiguration,
		NumberRequests:     numberRequests,
		Seconds:            seconds,
	}
}

func (h *Handler) HandlerRequests() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(h.chiMiddleware)
	r.Get("/", get)
	r.Post("/", post)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
}
