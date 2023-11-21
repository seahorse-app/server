package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"seahorse.app/server/database/models"
)

type JWTClaims struct {
	Sid string `json:"sid"`
	jwt.RegisteredClaims
}

// TODO: Proper error handling is needed asap

func AuthGuard(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session")

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// decrypt jwt token from cookie
		// check if session exists in database
		token, err := jwt.ParseWithClaims(cookie, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)

		if !ok {
			c.JSON(400, gin.H{"error": "Invalid jwt token"})
			c.Abort()
			return
		}

		// TODO: further and more detailed checks

		if claims.Sid == "" {
			c.JSON(400, gin.H{"error": "jwt token contains no session"})
			c.Abort()
			return
		}

		// get session from database
		// check if session is valid
		var session models.Session

		if err := DB.Where("sid=?", claims.Sid).First(&session).Error; err != nil {
			c.JSON(404, gin.H{"error": "Session not found"})
			c.Abort()
			return
		}

		c.Set("session", session.ID)
		c.Set("user", session.User)
	}
}
