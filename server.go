package main

import (
	"github.com/gin-gonic/gin"
	"seahorse.app/server/database"
	"seahorse.app/server/handlers"
	"seahorse.app/server/middleware"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// database setup
	database.SetupDatabase()

	r.GET("/middlewaretest", middleware.AuthGuard(database.DB), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "middlewaretest",
		})
	})

	userHandler := handlers.UserHandler{DB: database.DB}

	userGroup := r.Group("/user")
	{
		userGroup.POST("/create", userHandler.Create)
		userGroup.POST("/login", userHandler.Login)
		userGroup.GET("/profile", userHandler.Profile)
		userGroup.PATCH("/profile", userHandler.UpdateProfile)
	}

	r.Run()
}
