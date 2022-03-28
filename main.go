package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"main.go/api"
	"main.go/database"
)

func main(){
	database.Connect()
	defer database.Conn.Close()
	app := fiber.New()

	api.InitRoutes(app)

	log.Println(time.Now())

	log.Fatalln(app.Listen(":8080"))
}
