package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jeanSagaz/rate-limiter/internal/application/dto"
	"github.com/jeanSagaz/rate-limiter/internal/domain"
	"github.com/jeanSagaz/rate-limiter/pkg/infra/redis"
)

type Handler struct {
	Redis              *redis.RedisConnection
	TokenConfiguration []dto.TokenConfiguration
	NumberRequests     int
	Seconds            int
}

func NewHandler(redis *redis.RedisConnection,
	tokenConfiguration []dto.TokenConfiguration,
	numberRequests int,
	seconds int) *Handler {
	return &Handler{
		Redis:              redis,
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
	//w.Header().Add("Content-Type", "application/json")
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

func verifyIp(h *Handler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := r.RemoteAddr
	l := len(ip) - strings.LastIndex(ip, ":")
	ip = ip[:len(ip)-l]

	rdb, _ := h.Redis.Connect(ctx)
	value, _ := h.Redis.Get(rdb, ip)
	if len(value) == 0 {

		d := domain.Entity{Key: ip, Count: 1, Time: time.Now()}
		h.Redis.Set(rdb, ip, d, time.Second*time.Duration(h.Seconds))

	} else {

		data := domain.Entity{}
		json.Unmarshal([]byte(value), &data)

		data.Count = data.Count + 1
		h.Redis.Set(rdb, data.Key, data, time.Second*time.Duration(h.Seconds))

		// n := time.Now()
		// a := data.Time.Add((time.Second * time.Duration(h.Seconds)) + 1)
		// if n.After(a) && data.Count > h.NumberRequests {
		if data.Count > h.NumberRequests {
			w.WriteHeader(http.StatusTooManyRequests)
		}

	}
}

func verifyToken(h *Handler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	header := strings.TrimSpace(r.Header.Get("API_KEY"))

	// token
	token := dto.TokenConfiguration{}
	var tokens = h.TokenConfiguration
	for _, t := range tokens {
		if t.Token == header {
			token = t
			break
		}
	}

	rdb, _ := h.Redis.Connect(ctx)
	value, _ := h.Redis.Get(rdb, token.Token)
	if len(value) == 0 {

		e := domain.Entity{Key: token.Token, Count: 1, Time: time.Now()}
		h.Redis.Set(rdb, token.Token, e, time.Second*token.Seconds)

	} else {

		data := domain.Entity{}
		json.Unmarshal([]byte(value), &data)

		data.Count = data.Count + 1
		h.Redis.Set(rdb, data.Key, data, time.Second*token.Seconds)

		// n := time.Now()
		// a := data.Time.Add(time.Second * time.Duration(token.Seconds))
		// if n.After(a) && data.Count > token.NumberRequests {
		if data.Count > token.NumberRequests {
			w.WriteHeader(http.StatusTooManyRequests)
		}

	}
}

func rateLimiter(h *Handler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	header := strings.TrimSpace(r.Header.Get("API_KEY"))
	h.Redis.Context = ctx

	if header != "" {
		// token
		verifyToken(h, w, r)
	} else {

		// ip
		verifyIp(h, w, r)
	}
}

func (h *Handler) chiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rateLimiter(h, w, r)

		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
