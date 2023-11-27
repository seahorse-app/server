package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"seahorse.app/server/database"
	"seahorse.app/server/handlers"
	"seahorse.app/server/middleware"
)

func main() {
	app := fiber.New()

	database.SetupDatabase()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
		AllowOriginsFunc: nil,
	}))

	app.Use("/user/profile", middleware.AuthGuard())

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "pong",
		})
	})

	userHandler := handlers.UserHandler{DB: database.DB}

	userGroup := app.Group("/user")
	userGroup.Post("/create", userHandler.Create)
	userGroup.Post("/login", userHandler.Login)
	userGroup.Get("/profile", userHandler.Profile)

	app.Listen(":3000")
}
