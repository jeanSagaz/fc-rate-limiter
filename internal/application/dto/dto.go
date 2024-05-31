package dto

import "time"

type TokenConfiguration struct {
	Token          string        `json:"token"`
	NumberRequests int           `json:"numberRequests"`
	Seconds        time.Duration `json:"seconds"`
}
