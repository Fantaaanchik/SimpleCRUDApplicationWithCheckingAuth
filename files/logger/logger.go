package logger

import (
	"github.com/natefinch/lumberjack"
	"log"
)

// InitLogger - запускает логирование в файле logs.txt/ starts logging errors in the logs.txt file
func InitLogger() {
	filename := "files/logger/logs.txt"
	log.SetOutput(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    20,   // megabyte
		MaxBackups: 90,   //quantity of files
		MaxAge:     365,  //days
		Compress:   true, // disabled by default
	})
}
