package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"profOrientation/config"
	"profOrientation/models"
)

var db *gorm.DB

// ConnToDB - реализовано подключение к базе данных
func ConnToDB() *gorm.DB {
	var err error
	db, err = gorm.Open(postgres.Open(config.Conf.DB), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to DB")
	}
	autoMigrate()
	return db
}

// GetDb - получение базы данных и функционал доступный в нем
func GetDb() *gorm.DB {
	return db
}

// CloseDB - закрывает базу данных
func CloseDB() {
	if db == nil {
		return
	}
	PgDB, err := db.DB()
	err = PgDB.Close()
	if err != nil {
		log.Println("не удалось закрыть DB: ", err.Error())
	}
}

// autoMigrate - автоматическая миграция моделек в таблицу
func autoMigrate() {
	for _, model := range []interface{}{
		(*models.Answer)(nil),
		(*models.Profession)(nil),
		(*models.Question)(nil),
		(*models.Quiz)(nil),
		(*models.AuthInput)(nil),
		(*models.User)(nil),
		(*models.UserAnswer)(nil),
	} {
		err := db.AutoMigrate(model)
		if err != nil {
			log.Printf("Failed to migrate model %s: %s\n", model, err)
			return
		}
	}
	return
}
