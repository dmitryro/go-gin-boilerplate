// src/internal/models/role.go
package models

import (
    "time"
    "github.com/lib/pq"
)

// Role defines user roles and permissions.
// swagger:model
type Role struct {
    ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
    Name        string         `json:"name" gorm:"not null" example:"admin"`
    Permissions pq.StringArray `json:"permissions" gorm:"type:text[]" swaggertype:"array,string" example:"[\"read\",\"write\",\"delete\"]"`
    CreatedAt   time.Time      `json:"created_at" example:"2023-04-01T12:00:00Z"`
}
