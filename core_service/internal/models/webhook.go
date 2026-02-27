package models

type WebhookMessage struct {
	Id          string  `json:"id"`
	Type        string  `json:"type"`
	Timestamp   int64   `json:"timestamp"`
	Status      *string `json:"status,omitempty"`
	Description *string `json:"description,omitempty"`
}
