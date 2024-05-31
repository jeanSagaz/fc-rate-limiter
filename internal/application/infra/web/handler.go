package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jeanSagaz/rate-limiter/internal/application/dto"
	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis"
)

type Handler struct {
	IRedisRepository   redis.IRedisRepository
	TokenConfiguration []dto.TokenConfiguration
	NumberRequests     int
	Seconds            int
}

func NewHandler(iRedisRepository redis.IRedisRepository,
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

	response := dto.Response{Message: "Hello World"}
	json.NewEncoder(w).Encode(response)
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request dto.Request

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}
