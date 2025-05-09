package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/cart"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/license"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// PayForCart(cardId) —> заглушка
func PostAddToCart(service cart.Service, licenseService license.Service) fiber.Handler{
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil{
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err)) 
		}
		licenseId, err := c.ParamsInt("licenseId")
		if err != nil{
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(err)) 
		}

		//check if license belongs to beatmaker
		license, err := licenseService.ReadLicenseById(uint(licenseId))
		if err != nil{
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(err)) 
		}

		if userId == license.UserID {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(errors.New(
				"you cannot add your own beat to your cart"))) 
		}
		
		err = service.CreateProductToCart(userId, uint(licenseId))
		if err != nil{
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err)) 
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message":"successfully added to cart",
		}) 
	}
}

func GetCartByUser(service cart.Service) fiber.Handler{
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil{
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err)) 
		}
		cart, err := service.ReadCartByUser(userId)
		if err != nil{
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err)) 
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateSuccessResponse(cart)) 
	}
}

func DeleteLicenseFromCart(service cart.Service) fiber.Handler{
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil{
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err)) 
		}
		licenseId, err := c.ParamsInt("licenseId")
		if err != nil{
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(err)) 
		}
		err = service.DeleteFromCart(userId, uint(licenseId))
		if err != nil{
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err)) 
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message":"successfully deleted from cart",
		}) 
	}
}

func Hello(service cart.Service) fiber.Handler{
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message":"hello",
		}) 
	}
}


func getIdFromJWT(c *fiber.Ctx) (uuid.UUID, error){
	auth := c.GetReqHeaders()
	authHeader, ok := auth["Authorization"]
	if !ok {
		return uuid.Nil, errors.New("auth header is absent")
	}
	splitToken := strings.Split(authHeader[0], "Bearer ")
	tokenStr := splitToken[1]

	nilUuid := uuid.Nil
	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nilUuid, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nilUuid, err
	}

	id := claims["id"].(string)
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nilUuid, err
	}

	return uuid, nil
}



