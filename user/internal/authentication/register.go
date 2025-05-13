package authentication

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/user/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Email    string `json:"email" example:"eugene@example.com"`
	Password string `json:"password" example:"securepassword123"`
}

// Register godoc
// @Summary User Register
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body RegisterInput true "User credentials"
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 400 {object} presenters.UserErrorResponse
// @Failure 401 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /register [post]
func Register(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := RegisterInput{}

		err := c.BodyParser(&input)
		if err != nil {
			c.Status(http.StatusUnprocessableEntity)
			return c.JSON(presenters.CreateUserErrorResponse(err))
		}

		result, err := service.FetchUserByEmail(input.Email)

		if (err != nil) && !(errors.Is(err, gorm.ErrRecordNotFound)) {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateUserErrorResponse(err))
		}

		if result != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User already exists"})
		}

		user := entities.User{
			Email:    input.Email,
			Password: input.Password,
			RoleID:   1,
		}
		newUser, err := service.InsertUser(&user)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateUserErrorResponse(err))
		}

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = newUser.ID
		claims["role"] = newUser.ID
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

		return c.JSON(fiber.Map{"message": "Successfully registered", "token": t})
	}
}
