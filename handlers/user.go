package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"seahorse.app/server/database/models"
	"seahorse.app/server/utils"
)

// TODO: add validation for email

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

type UserProfileDTO struct {
	Birthdate string `json:"birthdate"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserProfileOwnerDTO struct {
	UserProfileDTO
	ID uuid.UUID `json:"id"`
}

func (handler *UserHandler) Create(c *fiber.Ctx) error {
	var userData UserCreateDTO
	if err := c.BodyParser(&userData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if userFound := handler.DB.Where("email=?", userData.Email).First(&models.User{}); userFound.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "User already exists",
		})
	}

	passwordHash, err := utils.HashPassword(userData.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	user := models.User{
		Email:        userData.Email,
		PasswordHash: passwordHash,
	}

	handler.DB.Create(&user)

	// TODO: send welcome mail to user
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
		"user":    user,
	})
}

func (handler *UserHandler) Login(c *fiber.Ctx) error {
	// TODO: replace with env variable for domain
	// TODO: check which device is logging in for longer/shorter session
	// TODO: set cookie expiration accourdingly

	var userLoginData UserLogin
	if err := c.BodyParser(&userLoginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	var user models.User
	if err := handler.DB.Where("email=?", userLoginData.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if !utils.CheckPassword(userLoginData.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	session := models.Session{
		UserID: user.ID,
	}

	handler.DB.Create(&session)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"sid": session.ID,
		"iss": "seahorse.app",
		"aud": "user",
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		handler.DB.Delete(&session)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	cookie.HTTPOnly = true
	cookie.Secure = false
	c.Cookie(cookie)
	return c.JSON(fiber.Map{
		"message": "Logged in",
	})
}

// TODO: corporate both profile functions into one

func (handler *UserHandler) Profile(c *fiber.Ctx) error {
	// TODO: check for authorization
	var user models.User
	// TODO: sanitize input
	userParam := c.Params("id")
	if err := handler.DB.Where("id=?", userParam).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": UserProfileDTO{
			Birthdate: user.BirthDate,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	})
}

//func (handler *UserHandler) OwnProfile(c *gin.Context) error {
//	user := c.MustGet("user").(models.User)
//	c.JSON(200, gin.H{"user": UserProfileOwnerDTO{
//		UserProfileDTO: UserProfileDTO{
//			Birthdate: user.BirthDate,
//			FirstName: user.FirstName,
//			LastName:  user.LastName,
//			Email:     user.Email,
//		},
//		ID: user.ID}})
//}

//func (handler *UserHandler) UpdateProfile(c *gin.Context) {
//	// TODO: check also if user is admin => then update other user than just him
////	user := c.MustGet("user").(models.User)
//	if err := handler.DB.Where("id=?", user.ID).First(&user).Error; err != nil {
//		c.JSON(404, gin.H{"error": "User not found"})
//		return
//	}
//	var UserProfileUpdateData UserProfileDTO
//	if err := c.BodyParser(&UserProfileUpdateData); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}

//if (UserProfileUpdateData == UserProfileDTO{}) {
//	c.JSON(400, gin.H{"error": "No data provided"})
//}

// TODO: more elegant way to do this
// TODO: if email is changed send confirmation mail
// TODO: send email to old mail to verify change

//	if user.BirthDate != UserProfileUpdateData.Birthdate && UserProfileUpdateData.Birthdate != "" {
//		user.BirthDate = UserProfileUpdateData.Birthdate
////	}
//
//	if user.FirstName != UserProfileUpdateData.FirstName && UserProfileUpdateData.FirstName != "" {
//		user.FirstName = UserProfileUpdateData.FirstName
//	}

//	if user.LastName != UserProfileUpdateData.LastName && UserProfileUpdateData.LastName != "" {
//		user.LastName = UserProfileUpdateData.LastName
//	}

//	if user.Email != UserProfileUpdateData.Email && UserProfileUpdateData.Email != "" {
//		user.Email = UserProfileUpdateData.Email
//	}

//	handler.DB.Save(&user)

//}
