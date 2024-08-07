package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitCors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // provide the to-be only accessible ip:port's or frontend host:port.
		AllowHeaders:     "Content-Type",
		AllowMethods:     "GET",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           3600,
	}))
}
