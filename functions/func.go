package functions

import (
	"fmt"
	"log"
	"profOrientation/db"
	"profOrientation/models"
)

// ProjectFunc - реалзация функционала проекта и вывод результата на терминал
func ProjectFunc() {
	defer db.CloseDB()
	var err error
	// Получение списка всех квизов
	var quizzes []models.Quiz
	fmt.Println("model created")

	err = db.GetDb().Raw("SELECT id, title FROM quizzes").Scan(&quizzes).Error
	if err != nil {
		log.Println("Failed to select id, title from quizzes, ERROR: ", err)
		return
	}
	totalScore := 0

	// Запуск квиза
	quizIndex := 0
	for {
		// Выбор текущего квиза
		QUIZ := quizzes[quizIndex]
		fmt.Println(QUIZ.Title)

		// Получение списка вопросов для текущего квиза
		var questions []models.Question
		err = db.GetDb().Raw("SELECT id, text FROM questions WHERE quiz_id = ?", QUIZ.ID).Scan(&questions).Error
		if err != nil {
			log.Println("Failed select to bring an id, text from model Question")
			return
		}

		// Получение ответов пользователя на каждый вопрос
		var answers []models.Answer
		for _, question := range questions {
			fmt.Print(question.ID)
			fmt.Print(". ")
			fmt.Println(question.Text)

			// Получение списка ответов для текущего вопроса
			var answerModels []models.Answer
			err = db.GetDb().Raw("SELECT id, text, score FROM answers WHERE question_id = ?", question.ID).Scan(&answerModels).Error
			if err != nil {
				log.Println("Failed select to bring an id, score from model Answer")
				return
			}
			for _, a := range answerModels {
				fmt.Printf("Варианты: %s\n", a.Text)
			}

			var answerID int
			_, err = fmt.Scan(&answerID)
			if err != nil {
				log.Println("Empty text area or unCorrect input from Scan at 59 row code")
				fmt.Println("Empty text area or unCorrect input from Scan at 60 row code")
				return
			}
			fmt.Println("answerID", answerID)
			// Получение выбранного ответа
			var a models.Answer
			err = db.GetDb().Raw("SELECT * FROM answers WHERE id = ?", answerID).Scan(&a).Error
			if err != nil {
				log.Println("Failed select to bring an answer score from model Answer")
				return
			}
			fmt.Println("a -", a.ID)
			if a.ID == 0 {
				log.Println("Selected answer doesn't exist")
				fmt.Println("Selected answer doesn't exist")
				return
			}

			answers = append(answers, a)

			totalScore += a.Score
		}

		// Получение списка профессий с баллами больше или равными суммарному баллу пользователя
		var professions []models.Profession
		err = db.GetDb().Raw("SELECT id, title, description, sum_score FROM professions WHERE sum_score <= ?", totalScore).Scan(&professions).Error
		if err != nil {
			log.Println("Failed select to bring an id, title, description, min_score from model Profession")
			return
		}
		// Вывод результатов квиза
		fmt.Println("Результаты квиза:")
		fmt.Printf("Вы получили %d баллов\n", totalScore)
		if len(professions) == 0 {
			fmt.Println("К сожалению, нет профессий, соответствующих вашим ответам :(")
		} else {
			fmt.Println("Подходящие профессии:")
			for _, profession := range professions {
				fmt.Printf("%s (%s)\n", profession.Title, profession.Description)
			}
		}
		// Запрос на повторение квиза
		fmt.Println("Хотите повторить квиз?")
		var again string
		_, err = fmt.Scan(&again)
		if err != nil {
			log.Println("Empty text area from Scan at 182 row code")
			fmt.Println("Empty text area from Scan at 182 row code")
			return
		}
		if again != "y" && again != "Y" && again != "н" && again != "Н" {
			break
		}
		// Выбор следующего квиза
		quizIndex = (quizIndex + 1) % len(quizzes)
	}
}
