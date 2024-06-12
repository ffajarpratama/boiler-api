package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID    uuid.UUID      `json:"user_id" gorm:"primaryKey; default:gen_random_uud()"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"-" gorm:"column:password"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`

	// json fields
	AccessToken  string `json:"access_token,omitempty" gorm:"-"`
	RefreshToken string `json:"refresh_token,omitempty" gorm:"-"`
}

func (User) TableName() string {
	return "tr_user"
}
