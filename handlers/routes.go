package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"profOrientation/config"
)

// StartRoutes - запускает роуты и проверяет, что они запустились правильно, также группирует роуты и вызывает функцию
// QuizAuthCheckMiddleware
func StartRoutes() {
	r := gin.Default()
	r.GET("/ping", ping)

	r.POST("/registration", Registration)
	r.POST("/login", Login)

	quizGroup := r.Group("/v1")
	quizGroup.Use(QuizAuthCheckMiddleware())
	{
		quizGroup.GET("/quiz", getQuiz)
		quizGroup.GET("/question", getQuestions)
		quizGroup.GET("/answers", getAnswers)
		quizGroup.GET("/total_score", getTotalScore)
		quizGroup.GET("/file", getRecommendedProfessionFile)
		quizGroup.POST("/user_answer", saveAnswerForEachQuestion)
		quizGroup.POST("/quiz", createNewQuiz)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "page not found"})
	})

	err := r.Run(config.Conf.App.PortRun)
	if err != nil {
		log.Fatal("router failed to start")
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Connection established!")
}
