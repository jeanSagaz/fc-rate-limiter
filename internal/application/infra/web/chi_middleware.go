package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jeanSagaz/rate-limiter/internal/application/dto"
	"github.com/jeanSagaz/rate-limiter/internal/domain"
)

func verifyIp(h *Handler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := r.RemoteAddr
	l := len(ip) - strings.LastIndex(ip, ":")
	ip = ip[:len(ip)-l]

	value, _ := h.IRedisRepository.Get(ctx, ip)
	if len(value) == 0 {

		d := domain.Entity{Key: ip, Count: 1, Time: time.Now()}
		h.IRedisRepository.Set(ctx, ip, d, time.Second*time.Duration(h.Seconds))

	} else {

		data := domain.Entity{}
		json.Unmarshal([]byte(value), &data)

		data.Count = data.Count + 1
		h.IRedisRepository.Set(ctx, data.Key, data, time.Second*time.Duration(h.Seconds))

		// n := time.Now()
		// a := data.Time.Add((time.Second * time.Duration(h.Seconds)) + 1)
		// if n.After(a) && data.Count > h.NumberRequests {
		if data.Count > h.NumberRequests {
			// w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			w.WriteHeader(http.StatusTooManyRequests)
			r.Body = io.NopCloser(bytes.NewReader([]byte("")))
		}

	}
}

func verifyToken(h *Handler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	header := strings.TrimSpace(r.Header.Get("API_KEY"))

	// token
	token := dto.TokenConfiguration{Token: "5753bac1-fad8-490a-818e-f97074655028", NumberRequests: 5, Seconds: 5}
	var tokens = h.TokenConfiguration
	for _, t := range tokens {
		if t.Token == header {
			token = t
			break
		}
	}

	value, _ := h.IRedisRepository.Get(ctx, token.Token)
	if len(value) == 0 {

		e := domain.Entity{Key: token.Token, Count: 1, Time: time.Now()}
		h.IRedisRepository.Set(ctx, token.Token, e, time.Second*token.Seconds)

	} else {

		data := domain.Entity{}
		json.Unmarshal([]byte(value), &data)

		data.Count = data.Count + 1
		h.IRedisRepository.Set(ctx, data.Key, data, time.Second*token.Seconds)

		// n := time.Now()
		// a := data.Time.Add(time.Second * time.Duration(token.Seconds))
		// if n.After(a) && data.Count > token.NumberRequests {
		if data.Count > token.NumberRequests {
			// w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			w.WriteHeader(http.StatusTooManyRequests)
			r.Body = io.NopCloser(bytes.NewReader([]byte("")))
		}

	}
}

func rateLimiter(h *Handler, w http.ResponseWriter, r *http.Request) {
	header := strings.TrimSpace(r.Header.Get("API_KEY"))

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
