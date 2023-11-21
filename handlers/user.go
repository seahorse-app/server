package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

type UserCreate struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (handler *UserHandler) Create(c *gin.Context) {
	var userData UserCreate
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	println(userData.Username)
}
