package main

import (
	"profOrientation/config"
	"profOrientation/db"
	"profOrientation/files/logger"
	"profOrientation/functions"
	"profOrientation/handlers"
)

func main() {
	config.ReadConf("config/config.json")
	logger.InitLogger()
	db.ConnToDB()
	handlers.StartRoutes()
	functions.ProjectFunc()
}
