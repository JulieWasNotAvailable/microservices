package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/entities"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

// HandleGoogleAuth godoc
// @Summary Authenticate with Google
// @Description Authenticate user using Google OAuth token and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body presenters.Data true "Google OAuth token"
// @Success 200 {object} map[string]interface{} "Returns JWT token and user info"
// @Success 201 {object} map[string]interface{} "Returns JWT token and user info (new user)"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 401 {object} map[string]interface{} "Invalid Google token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/google/getjwt [post]
func HandleGoogleAuth(service user.Service) fiber.Handler {
	return func (c *fiber.Ctx) error {

	var json presenters.Data
	if err := c.BodyParser(&json); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Verify Google token
	payload, err := idtoken.Validate(context.Background(), json.Token, os.Getenv("Client"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Google token"})
	}

	// Extract user data
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	userId, _ := payload.Claims["sub"].(string)

	user, err := service.FetchUserByEmail(email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(presenters.CreateUserErrorResponse(err))
	}

	if user == nil {
		newUser := entities.User{
			Email : email,
		}
		result, err := service.InsertUser(&newUser)
		if err != nil{
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateUserErrorResponse(err))
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id": result.ID,
			"role": result.RoleID,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		})
		
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}
	
		return c.JSON(fiber.Map{
			"token": tokenString,
			"user": fiber.Map{
				"id":    userId,
				"email": email,
				"name":  name,
			},
		})
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"role": user.RoleID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":    userId,
			"email": email,
			"name":  name,
		},
	})
}
}