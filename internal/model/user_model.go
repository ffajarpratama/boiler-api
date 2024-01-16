package model

import (
	"github.com/google/uuid"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	UserID    uuid.UUID             `json:"user_id" gorm:"primaryKey; default:gen_random_uuid()"`
	Name      string                `json:"name"`
	Email     string                `json:"email"`
	Password  string                `json:"-"`
	CreatedAt int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int                   `json:"updated_at" gorm:"autoCreateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-" gorm:"column:deleted_at"`

	// json fields
	Token string `json:"token,omitempty" gorm:"-"`
}

func (User) TableName() string {
	return "tm_user"
}
