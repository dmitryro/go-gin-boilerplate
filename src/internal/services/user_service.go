package services

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"go_api/internal/models"
)

type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new instance of UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// PreloadRole preloads the Role association for a given user
func (s *UserService) PreloadRole(user *models.User) error {
	return s.db.Model(user).Association("Role").Find(&user.Role)
}

// GetAllUsers retrieves all users from the database
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := s.db.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser creates a new user in the database
func (s *UserService) CreateUser(userCreateRequest *models.UserCreateRequest) (*models.User, error) {
	user := &models.User{
		Email:    userCreateRequest.Email,
		Username: userCreateRequest.Username,
		Password: userCreateRequest.Password, // Plaintext temporarily
		First:    userCreateRequest.First,
		Last:     userCreateRequest.Last,
		Phone:    userCreateRequest.Phone,
		RoleID:   userCreateRequest.RoleID,
	}

	plaintextPassword := user.Password

	log.Printf("CreateUser: Creating user with email: %s", user.Email)

	// Create user with plaintext password initially
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// Hash password and update
	if err := s.UpdateUserPassword(user.Username, plaintextPassword); err != nil {
		log.Printf("CreateUser: Failed to update password for user '%s': %v", user.Username, err)
		// Consider returning error to caller here if preferred
	}

	// Preload Role association
	if err := s.db.Model(user).Association("Role").Find(&user.Role); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by its ID from the database
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}

// DeleteUser deletes a user by its ID from the database
func (s *UserService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

// GetUserByEmail retrieves a user by its email from the database
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by its username from the database
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUsersByRoleID retrieves all users who have the specified role ID
func (s *UserService) GetUsersByRoleID(roleID uint) ([]models.User, error) {
	var users []models.User
	if err := s.db.Preload("Role").Where("role_id = ?", roleID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUserPassword updates the password for a user, given the username and the new password.
func (s *UserService) UpdateUserPassword(username, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.db.Model(&models.User{}).Where("username = ?", username).Update("password", string(hashedPassword)).Error
}

// ChangeUserPassword changes the password for a user, given the username, the old password, and the new password.
func (s *UserService) ChangeUserPassword(username, oldPassword, newPassword string) error {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	// Compare the provided old password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password
	user.Password = string(hashedPassword)
	return s.db.Save(&user).Error
}

