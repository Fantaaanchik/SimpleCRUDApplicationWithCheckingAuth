package models

import (
	"time"
)

type Question struct {
	ID        int        `gorm:"column:id" json:"id"`
	QuizID    int        `gorm:"column:quiz_id" json:"quiz_id"`
	Text      string     `gorm:"column:text" json:"text"`
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
}

func (q Question) TableName() string {
	return "questions"
}

type QuestionWithResp struct {
	ID      int    `gorm:"column:id" json:"id"`
	QuizID  int    `gorm:"column:quiz_id" json:"quiz_id"`
	Text    string `gorm:"column:text" json:"text"`
	Answers []Answer
}
