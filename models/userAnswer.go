package models

import "time"

type UserAnswer struct {
	ID         int        `gorm:"column:id" json:"id"`
	QuestionID int        `gorm:"column:question_id" json:"question_id"`
	Score      int        `gorm:"column:score" json:"score"`
	Option     int        `gorm:"column:option" json:"option"`
	CreatedAt  *time.Time `gorm:"autoCreateTime"`
	UpdatedAt  *time.Time `gorm:"autoUpdateTime"`
}

func (u UserAnswer) TableName() string {
	return "user_answers"
}
