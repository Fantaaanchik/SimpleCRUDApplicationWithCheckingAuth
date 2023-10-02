package models

import "time"

type AuthInput struct {
	Username  string     `json:"username" binding:"required"`
	Password  string     `json:"password" binding:"required"`
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
}

func (a AuthInput) TableName() string {
	return "auth_inputs"
}
