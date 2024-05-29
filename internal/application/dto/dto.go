package dto

import "time"

type Request struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Response struct {
	Message string `json:"message"`
}

type TokenConfiguration struct {
	Token          string        `json:"token"`
	NumberRequests int           `json:"numberRequests"`
	Seconds        time.Duration `json:"seconds"`
}
