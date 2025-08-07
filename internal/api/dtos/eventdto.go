package api

import (
	"time"
)

type EventType string

type Source string

type Data struct {
	Action   string                 `json:"action"`
	Value    float32                `json:"value"`
	Metadata map[string]interface{} `json:"metadata"`
}

type EventDTO struct {
	ID        *string   `json:"id"`
	Type      EventType `json:"type"`
	Source    Source    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
	UserID    *string   `json:"user_id"`
	Data      Data      `json:"data"`
}
