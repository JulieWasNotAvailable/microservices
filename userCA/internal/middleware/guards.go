package middleware

import (
	"log"
	"os"
	// "strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// 	jwtware.New — это функция из пакета jwtware, которая создает middleware для проверки JWT.
// Она принимает конфигурацию (jwtware.Config) и возвращает middleware типа func(*fiber.Ctx) error.
// middleware будет выполняться перед обработкой запроса в защищенных маршрутах.
// Вовзращаемая функция достаёт JWT из запроса, проверяет его с использованием предоставленного ключа, а если что-то не так вызывает error handler

func LoadENV() error {
	err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	return err
}

func ProtectedRequiresBeatmaker() func(*fiber.Ctx) error {
	LoadENV()
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
	LoadENV()
	log.Println("inside protected req beatmaker")
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

// func ProtectedGoogle() func(*fiber.Ctx) error {
// 	return func(c *fiber.Ctx) error {
// 	authHeader := c.Get("Authotization")
// 	if authHeader == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//             "error": "Authorization header missing",
//         })
// 	}

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//     // if tokenString == authHeader {
//     //     return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//     //         "error": "Invalid Authorization header format",
//     //     })
//     // }
	
// 	oauth2Service, err := oauth2.NewService(c.Context(), option.WithoutAuthentication())
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "error": "Failed to create OAuth2 service",
//         })
//     }

// 	tokenInfo, err := oauth2Service.Tokeninfo().AccessToken(tokenString).Do()
//     if err != nil {
//         return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//             "error": "Invalid access token",
//             "details": err.Error(),
//         })
//     }
	
// 	if tokenInfo.Audience != "your-client-id.apps.googleusercontent.com" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Token not issued for this application",
// 		})
// 	}

// 	// Check token expiration
// 	if tokenInfo.ExpiresIn <= 0 {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Token expired",
// 		})
// 	}

// 	// Verify email if needed
// 	if !tokenInfo.EmailVerified {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Email not verified",
// 		})
// 	}
// }
// }