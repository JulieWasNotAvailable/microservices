package guards

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
		return c.JSON(fiber.Map{"status": "false", "message": "Missing or malformed JWT", "data": nil})
	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "false", "message": "Invalid or expired JWT", "data": nil})
	}
}

// func CustomJWTMiddleware() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		LoadENV()
// 		secret := os.Getenv("SECRET")
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Missing Authorization header",
// 			})
// 		}

// 		tokenString := authHeader[len("Bearer "):]

// 		log.Println(tokenString)
// 		parts := strings.Split(tokenString, ".")
// 		signature, err := base64.RawURLEncoding.DecodeString(parts[2])
// 		if err != nil {
// 			return err
// 		}
// 		myString := string(signature[:])
// 		log.Println(myString)

// 		// return jwtware.New(jwtware.Config{
// 		// 	SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("SECRET"))},
// 		// 	ErrorHandler: jwtError,
// 		// })
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, errors.New("unexpected signing method")
// 			}
// 			return []byte(secret), nil
// 		})

// 		if err != nil {
// 			return err
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 			for key, val := range claims {
// 				if f, ok := val.(float64); ok {
// 					// Check if it's a whole number before converting
// 					if f == float64(int64(f)) {
// 						claims[key] = int64(f)
// 					}
// 				}
// 			}

// 			if exp, ok := claims["exp"].(int64); ok {
// 				if time.Now().Unix() > exp {
// 					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 						"error": "Token expired",
// 					})
// 				}
// 			} else {
// 				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 					"error": "Invalid exp claim",
// 				})
// 			}
// 		} else {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid token claims",
// 			})
// 		}

// 		// 5. Store token in locals and proceed
// 		c.Locals("user", token)
// 		return c.Next()
// 	}
// }
