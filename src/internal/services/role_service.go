package services

import (
	"errors"

	"gorm.io/gorm"
	"go_api/internal/models"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

// GetAllRoles retrieves all roles from the database
func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID retrieves a role by its ID from the database
func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetRoleByName retrieves a role by its name from the database
func (s *RoleService) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	if err := s.db.Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// CreateRole creates a new role in the database
func (s *RoleService) CreateRole(role *models.Role) error {
	return s.db.Create(role).Error
}

