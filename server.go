package main

import (
	"github.com/gin-gonic/gin"
	"seahorse.app/server/database"
	"seahorse.app/server/handlers"
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
	}

	r.Run()
}
