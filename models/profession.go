package models

import "time"

type Profession struct {
	ID          int        `gorm:"column:id" json:"id"`
	Title       string     `gorm:"column:title" json:"title"`
	Description string     `gorm:"column:description" json:"description"`
	Skills      string     `gorm:"skills" json:"skills"`
	SumScore    int        `gorm:"sum_score" json:"sum_score"`
	CreatedAt   *time.Time `gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime"`
}

func (p Profession) TableName() string {
	return "professions"
}
