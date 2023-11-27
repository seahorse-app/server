package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"seahorse.app/server/database"
	"seahorse.app/server/database/models"
)

type JWTClaims struct {
	ID uuid.UUID `json:"sid"`
	jwt.RegisteredClaims
}

// TODO: Proper error handling is needed asap

func AuthGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")
		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		token, err := jwt.ParseWithClaims(cookie, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		claims, ok := token.Claims.(*JWTClaims)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		var session models.Session

		if err := database.DB.Where("id=?", claims.ID).First(&session).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		var user models.User

		if err := database.DB.Preload("Sessions").Where("id=?", session.UserID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		c.Locals("user", user)
		c.Locals("session", session)

		return c.Next()
	}
}
