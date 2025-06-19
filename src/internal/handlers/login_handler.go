// src/internal/handlers/login_handler.go

package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go_api/internal/models"
	"go_api/internal/services"
)

// LoginHandler handles login requests
type LoginHandler struct {
	service *services.LoginService
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(service *services.LoginService) *LoginHandler {
	return &LoginHandler{service: service}
}

// Login godoc
// @Summary     Authenticate user and return JWT token
// @Description Validates user credentials and returns a JWT token upon successful authentication
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       credentials body models.LoginRequest true "User credentials"
// @Success     200 {object} models.TokenResponse
// @Failure     401 {object} models.ErrorResponse
// @Router      /api/login [post]
func (h *LoginHandler) Login(c *gin.Context) {
	var creds models.LoginRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Println("Error binding credentials:", err) // Log any binding errors
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid request"})
		return
	}

	// Attempt authentication with the service
	token, err := h.service.Authenticate(creds.Username, creds.Password)
	if err != nil {
		log.Println("Authentication failed for username:", creds.Username) // Log authentication failure
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Code: http.StatusUnauthorized, Message: "Invalid credentials"})
		return
	}

	log.Println("Authentication successful for username:", creds.Username) // Log successful authentication
	c.JSON(http.StatusOK, models.TokenResponse{Token: token})
}
