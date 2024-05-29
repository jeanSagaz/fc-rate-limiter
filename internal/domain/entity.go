package domain

import "time"

type Entity struct {
	Key   string    `json:"key"`
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}
