package dto

import (
	"time"

	"github.com/google/uuid"
)

type SendToWebhookDto struct {
	UID       uuid.UUID `json:"user_id"`
	IpAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
}
