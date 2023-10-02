package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"profOrientation/db"
	"profOrientation/models"
	"strings"
)

// Registration - регистрирует нового пользователя
func Registration(c *gin.Context) {
	var (
		authInput models.AuthInput
		u         models.User
		err       error
	)

	err = c.ShouldBindJSON(&authInput)
	if err != nil {
		log.Println("bind err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(authInput)

	u.Login = authInput.Username
	u.Password = authInput.Password
	u.Password, err = hashPassword(u.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = db.GetDb().Create(&u).Error
	if err != nil {
		log.Println("create new user db err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Новый пользователь зарегистрирован!"})
}

// hashPassword - эта функция хэширует пароли
func hashPassword(password string) (newPass string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println("hash err:", err.Error())
		return "", err
	}
	newPass = strings.TrimSpace(string(hash))
	return
}

// Login - вход существующего пользователя
func Login(c *gin.Context) {
	var (
		input        models.AuthInput
		existingUser models.User
		err          error
	)
	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Println("bind err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = db.GetDb().Where("login = ?", input.Username).Find(&existingUser).Error
	if err != nil {
		log.Println("find user login err:", err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input)
	err = checkingHashFromPasswordToken(existingUser.Password, input.Password)
	if err != nil {
		log.Println("hash check err:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Пароль неверный!"})
		return
	}
	token, err := GenerateToken(uint(existingUser.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Не удалось создать токен"})
		return
	}
	c.JSON(200, token)
}

// checkingHashFromPasswordToken - проверка хэша из токена пароля
func checkingHashFromPasswordToken(existing, input string) error {
	return bcrypt.CompareHashAndPassword([]byte(existing), []byte(input))
}
