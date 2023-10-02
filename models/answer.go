package models

import "time"

type Answer struct {
	ID         int        `gorm:"column:id" json:"id"`
	QuestionID int        `gorm:"column:question_id" json:"question_id"`
	Text       string     `gorm:"column:text" json:"text"`
	Score      int        `gorm:"column:score" json:"score"`
	CreatedAt  *time.Time `gorm:"autoCreateTime"`
	UpdatedAt  *time.Time `gorm:"autoUpdateTime"`
}

func (a Answer) TableName() string {
	return "answers"
}
