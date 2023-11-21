package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"seahorse.app/server/database/models"
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

	// TODO: Password encryption

	user := models.User{
		Email:        userData.Email,
		PasswordHash: userData.Password,
	}

	handler.DB.Create(&user)

	c.String(200, string(user.CreatedAt.String()))
}
