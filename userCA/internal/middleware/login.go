package middleware

import (
	"log"
	"os"
	"time"

	"github.com/JulieWasNotAvailable/microservices/user/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// @Description Login Input
type LoginInput struct {
	Email string `json:"email" example:"john_molly@example.com"`
	Password string `json:"password" example:"securepassword123"`
}

// Login get user and password
// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginInput true "User credentials"
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 400 {object} presenters.UserErrorResponse
// @Failure 401 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /login [post]
func Login(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
	input := LoginInput{}

	err := c.BodyParser(&input)

	if err != nil {
		return err
	}

	login := input.Email
	pass := input.Password

	user, err := service.FetchUserByEmail(login)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User not found"})
	}

	//to-do: use bcrypt
	if user.Password != pass {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid password"})
	}

	//creates an empty token having the signing method hs256
	// HS256 (HMAC with SHA-256) is a symmetric keyed hashing algorithm that uses one secret key
	token := jwt.New(jwt.SigningMethodHS256)

	//JWT consists of three parts: a header, payload, and signature
	//claim is one key-value pair in the payload
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.ID
	claims["role"] = user.RoleID
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	err = godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"message": "Success login", "token": t})
	}
}