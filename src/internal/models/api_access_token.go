package models

import (
    "time"
)

// ApiAccessToken represents API access credentials.
type ApiAccessToken struct {
    ID        uint      `json:"id"`
    Token     string    `json:"token"`
    UserID    uint      `json:"user_id"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}
