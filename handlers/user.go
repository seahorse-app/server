package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"seahorse.app/server/database/models"
	"seahorse.app/server/utils"
)

type UserHandler struct {
	DB *gorm.DB
}

type UserCreate struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (handler *UserHandler) Create(c *gin.Context) {
	var userData UserCreate
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if userFound := handler.DB.Where("email=?", userData.Email).First(&models.User{}); userFound.RowsAffected > 0 {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}

	passwordHash, err := utils.HashPassword(userData.Password)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:        userData.Email,
		PasswordHash: passwordHash,
	}

	handler.DB.Create(&user)

	// TODO: send welcome mail to user
	c.JSON(200, gin.H{"user": user})
}
