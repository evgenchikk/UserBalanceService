package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/evgenchikk/avito-tech-internship_backend_2022/app/handlers"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/config"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/db"
	"github.com/gin-gonic/gin"
)

func Run(conf *config.Config) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DB.DB_USER, conf.DB.DB_PASSWD, conf.DB.DB_HOST, conf.DB.DB_PORT, conf.DB.DB_NAME)
	dbConn, err := db.New(dbUrl)

	// app won't start if it failed to connect to db
	if err != nil {
		log.Println("Start failed")
		return
	}

	// try create a dir for storing reports
	_, err = os.Stat(conf.ReportDirPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(conf.ReportDirPath, 0755); err != nil {
			log.Println("Create report dir failed:", err.Error())
			conf.ReportDirPath = "."
		}
	}

	// try create a dir for storing logs
	_, err = os.Stat(conf.LogDirPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(conf.LogDirPath, 0755); err != nil {
			log.Println("Create log dir failed:", err.Error())
			conf.LogDirPath = "."
		}
	}

	// try create a log file
	logFile, err := os.Create(fmt.Sprintf("%s/%s.log", conf.LogDirPath, time.Now().Format("06-01-02_15:04:05")))
	if err == nil {
		gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	}
	defer logFile.Close()

	// bootstrap http server
	router := gin.Default()
	router.SetTrustedProxies(nil)
	handlers.RegisterRoutes(router, dbConn, conf.ReportDirPath, conf.Host+conf.Port)

	router.Run(conf.Port)
}
