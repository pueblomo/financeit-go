package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"main.go/api"
	"main.go/conf"
	"main.go/database"
	"main.go/sheduledTask"
)

func main(){
	conf.InitEnv()

	database.Connect()
	defer database.Conn.Close()
	app := fiber.New()

	api.InitRoutes(app)

	ticker := startSheduledTask()
	defer ticker.Stop()

	log.Fatalln(app.Listen(":8080"))
}

func startSheduledTask() *time.Ticker{
	ticker := time.NewTicker(8 * time.Hour)
	go func() {
    	for range ticker.C {
        	sheduledTask.SendToAggregator()
    	}
	}()
	return ticker
}
