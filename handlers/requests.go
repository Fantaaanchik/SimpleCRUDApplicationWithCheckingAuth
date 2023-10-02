package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"profOrientation/db"
	"profOrientation/models"
	"strconv"
)

// getQuiz - берет квизы, доступные в базе данных или показывает необходимый через id
func getQuiz(c *gin.Context) {
	var (
		quizzes []models.Quiz
		err     error
	)
	_, err = ExtractTokenByID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idFromQueryParam := c.Query("id")
	if idFromQueryParam == "" {
		err = db.GetDb().Find(&quizzes).Error
	} else {
		err = db.GetDb().Where("id = ?", idFromQueryParam).Find(&quizzes).Error
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if quizzes == nil {
		quizzes = []models.Quiz{}
	}
	c.JSON(http.StatusOK, gin.H{"quizzes": quizzes})
	return

}

// getQuestions -  берет вопросы, доступные в базе данных или фильтрует необходимые через id
func getQuestions(c *gin.Context) {
	var (
		question []models.Question
		err      error
	)
	_, err = ExtractTokenByID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idFromQuestionsFromQueryParam := c.Query("id")
	if idFromQuestionsFromQueryParam == "" {
		// Если параметр id не указан, получаем все вопросы
		err = db.GetDb().Find(&question).Error
	} else {
		// Если параметр id указан, получаем только вопрос с данным id
		err = db.GetDb().Where("id = ?", idFromQuestionsFromQueryParam).Find(&question).Error
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if len(question) == 0 {
		// Если вопрос не найден, отправляем соответствующий ответ
		c.JSON(http.StatusNotFound, gin.H{"message": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"question": question})
	return

}

// getAnswers -  берет ответы, доступные в базе данных или показывает необходимые варианты через question_id
func getAnswers(c *gin.Context) {
	var (
		answer []models.Answer
		err    error
	)
	_, err = ExtractTokenByID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idFromAnswersForEachQuestionFromQueryParam := c.Query("id")
	if idFromAnswersForEachQuestionFromQueryParam == "" {
		// Если параметр id не указан, получаем все варианты ответов
		err = db.GetDb().Find(&answer).Error
	} else {
		// Если параметр id указан, получаем только ответы с данным id
		err = db.GetDb().Select("id, question_id, text").Where("question_id = ?", idFromAnswersForEachQuestionFromQueryParam).Find(&answer).Error
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if len(answer) == 0 {
		// Если вопрос не найден, отправляем соответствующий ответ
		c.JSON(http.StatusNotFound, gin.H{"message": "Question not found"})
		return
	}
	// Исключаем поля CreatedAt, UpdatedAt и score из структуры ответа на время показа в запросе
	var answerResponses []map[string]interface{}
	for _, answer := range answer {
		answerResponse := map[string]interface{}{
			"id":          answer.ID,
			"question_id": answer.QuestionID,
			"text":        answer.Text,
		}
		answerResponses = append(answerResponses, answerResponse)
	}
	c.JSON(http.StatusOK, gin.H{"answer": answerResponses})
	return

}

// saveAnswerForEachQuestion - сохраняет ответ на каждый вопрос в таблицу user_answers
func saveAnswerForEachQuestion(c *gin.Context) {
	var input models.UserAnswer
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	// добавляем ответ пользователя в таблицу user_answers
	err = db.GetDb().Where("id = ?", input.ID).Create(&input).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	//Находим соответсвующий ответ из таблицы answers и обновляем оценку в таблице user_answers
	var answer models.Answer
	err = db.GetDb().Where("id = ?", input.ID).First(&answer).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	err = db.GetDb().Model(&input).Update("score", answer.Score).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ваш ответ записан, переходите к следующему вопросу"})
}

// createNewQuiz - создает новый квиз и записывает его в таблицу quizzes
func createNewQuiz(c *gin.Context) {
	var newQuiz models.Quiz
	if err := c.BindJSON(&newQuiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Записываем новый quiz в таблицу quizzes
	if err := db.GetDb().Create(&newQuiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quiz успешно записан в базу!"})
}

// getTotalScore - получает из таблицы user_answers общую сумму заработанных баллов и выводи его пользователю и
// рекомендует ему профессии исходя из полученных баллов
func getTotalScore(c *gin.Context) {
	var totalScore int
	err := db.GetDb().Table("user_answers").Select("sum(score)").Scan(&totalScore).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	var professions []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Skills      string `json:"skills"`
	}
	err = db.GetDb().Table("professions").Where("sum_score <= ?", totalScore).Order("sum_score desc").Select("title, description, skills").Find(&professions).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if len(professions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"total_score": totalScore, "message": "По набранным баллам подходяших профессий нет найдено!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"total_score": totalScore, "professions": professions})
	}
}

// getRecommendedProfessionFile - пользователь получает файл, название которого он ввел в GET запросе, с подробным описанием выбранной профессии
func getRecommendedProfessionFile(c *gin.Context) {
	// Получаем имя файла из параметров запроса
	fileName := c.Query("file")

	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File name is required"})
		return
	}

	// Составляем полный путь к файлу
	directory := "files/" + fileName

	// Открываем файл
	file, err := os.Open(directory)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()

	// Создаем срез байтов, в который запишем данные из файла для определения параметров файла
	bytes := make([]byte, 512)
	_, err = file.Read(bytes)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Указываем, что закреплен файл (вложение) в соответствующем header
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", fileName))

	// Определяем тип контента (нужен для определения формата файла)
	contentType := http.DetectContentType(bytes)
	c.Header("Content-Type", contentType)

	// Определяем размер файла
	FileStat, _ := file.Stat()                         // Информация о файле
	FileSize := strconv.FormatInt(FileStat.Size(), 10) // Размер
	// Указываем размер в соответствующем заголовке
	c.Header("Content-Length", FileSize)

	_, err = file.Seek(0, 0) // Прочтем данные файла еще раз
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Отправляем файл пользователю в качестве ответа на запрос
	c.File(directory)
}
