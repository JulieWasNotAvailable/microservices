package internal

import (
	"log"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func LoadENV() error {
	err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	return err
}

func ProtectedRequiresBeatmaker() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("SECRET"))},
		ErrorHandler: jwtError,
		SuccessHandler:  func(c *fiber.Ctx) error {
			
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			if claims["role"].(float64) != 2 {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                    "error": "Access denied: 'beatmaker' claim is missing or false",
                })
			}
			return c.Next()
		},
	})
}

func ProtectedRequiresModerator() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("SECRET"))},
		ErrorHandler: jwtError, 
		SuccessHandler:  func(c *fiber.Ctx) error {
			
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			if claims["role"].(float64) != 3 {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                    "error": "Access denied: 'admin' claim is missing or false",
                })
			}
			return c.Next()
		},
	})
}

func Protected() func(*fiber.Ctx) error {
	LoadENV()
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("SECRET"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}
