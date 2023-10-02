package models

import "time"

type User struct {
	ID        int64      `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	RoleID    int64      `json:"roleID" gorm:"column:role_id" validate:"gte=1"`
	FIO       string     `json:"fio" gorm:"column:fio"`
	Phone     string     `json:"phone" gorm:"column:phone"`
	Login     string     `json:"login" gorm:"column:login" validate:"gt=0"`
	Password  string     `json:"password" gorm:"column:password" validate:"gt=0"`
	Salt      string     `json:"salt" gorm:"column:salt"`
	Token     string     `json:"token" gorm:"column:token"`
	Active    bool       `json:"active" gorm:"column:active;default:true"`
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
}

func (t User) TableName() string {
	return "users"
}
