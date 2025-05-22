package handlers

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

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