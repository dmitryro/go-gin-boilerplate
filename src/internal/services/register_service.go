package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"go_api/internal/models"
	"gorm.io/gorm"
)

type RegisterService struct {
	db *gorm.DB
}

func NewRegisterService(db *gorm.DB) *RegisterService {
	return &RegisterService{db: db}
}

func (s *RegisterService) RegisterUser(req *models.UserCreateRequest) (*models.User, error) {
	// Check if username or email already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("username or email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		First:    req.First,
		Last:     req.Last,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		RoleID:   req.RoleID, // you may want to default to a role here, e.g. "guest" role id
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

