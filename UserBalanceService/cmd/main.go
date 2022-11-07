package main

import (
	"github.com/evgenchikk/avito-tech-internship_backend_2022/app"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/config"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/docs"

	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// @title           User Balance Service
// @version         1.0
// @description     A service that can perform some operations with user balances.

// @contact.name   Evgeny Belonogov
// @contact.email  ewbelonogov@ya.ru
// @contact.url    https://www.t.me/evgenchikkkkkk

// @BasePath /
// @schemes http
func main() {
	conf := config.New()

	docs.SwaggerInfo.Host = conf.Host + conf.Port

	app.Run(conf)
}
