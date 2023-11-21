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

type UserBaseDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserCreateDTO struct {
	UserBaseDTO
}

type UserLogin struct {
	UserBaseDTO
}

func (handler *UserHandler) Create(c *gin.Context) {
	var userData UserCreateDTO
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

func (handler *UserHandler) Login(ctx *gin.Context) {
	var userLoginData UserLogin
	if err := ctx.ShouldBindJSON(&userLoginData); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := handler.DB.Where("email=?", userLoginData.Email).First(&user).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPassword(userLoginData.Password, user.PasswordHash) {
		ctx.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	// TODO: generate jwt token..
}
