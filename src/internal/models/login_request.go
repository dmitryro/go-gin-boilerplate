// src/internal/models/login.go
package models

// LoginRequest represents the payload for user login
// swagger:model
type LoginRequest struct {
    Username string `json:"username" example:"admin"`
    Password string `json:"password" example:"password123"`
}
