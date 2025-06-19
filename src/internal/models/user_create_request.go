package models

// UserCreateRequest represents the payload for creating a new user
type UserCreateRequest struct {
    Email    string `json:"email" example:"user@example.com"`
    Username string `json:"username" example:"johndoe"`
    Password string `json:"password" example:"strongpassword123"`
    First    string `json:"first" example:"John"`
    Last     string `json:"last" example:"Doe"`
    Phone    string `json:"phone" example:"+1234567890"`
    RoleID   uint   `json:"role_id" example:"1"`
}
