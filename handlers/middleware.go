package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// QuizAuthCheckMiddleware - проверяет правильность введенных пользователем данных для авторизации
func QuizAuthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ValidateToken(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
