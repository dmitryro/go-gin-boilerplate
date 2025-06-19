package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go_api/internal/models"
	"go_api/internal/services"
)

// RegisterHandler handles user registration requests
type RegisterHandler struct {
	registerService *services.RegisterService
}

func NewRegisterHandler(registerService *services.RegisterService) *RegisterHandler {
	return &RegisterHandler{registerService: registerService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "User registration info"
// @Success 201 {object} map[string]interface{} "Created user"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/register [post]
func (h *RegisterHandler) Register(c *gin.Context) {
	var userCreateRequest models.UserCreateRequest
	if err := c.ShouldBindJSON(&userCreateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.registerService.RegisterUser(&userCreateRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"id":       createdUser.ID,
		"username": createdUser.Username,
		"email":    createdUser.Email,
		"first":    createdUser.First,
		"last":     createdUser.Last,
		"phone":    createdUser.Phone,
		"role_id":  createdUser.RoleID,
	}

	c.JSON(http.StatusCreated, response)
}

