package main

import (
	"github.com/gin-gonic/gin"
	"seahorse.app/server/database"
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

}
