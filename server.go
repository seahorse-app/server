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

	userHandler := handlers.UserHandler{DB: database.DB}

	userGroup := r.Group("/user")
	{
		userGroup.POST("/create", userHandler.Create)
		userGroup.POST("/login", userHandler.Login)
		// TODO: let authguard get database instance by itself
		userGroup.GET("/profile", middleware.AuthGuard(database.DB), userHandler.OwnProfile)
		userGroup.GET("/profile/:id", middleware.AuthGuard(database.DB), userHandler.Profile)
		userGroup.PATCH("/profile", middleware.AuthGuard(database.DB), userHandler.UpdateProfile)
	}

	r.Run()
}
