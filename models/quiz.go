package models

import "time"

type Quiz struct {
	ID        int        `gorm:"column:id" json:"id"`
	Title     string     `gorm:"column:title" json:"title"`
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
}

func (q Quiz) TableName() string {
	return "quizzes"
}
