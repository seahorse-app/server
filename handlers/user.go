package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (handler *UserHandler) Create(c *gin.Context) {
	c.String(200, "Hello World!")
}
