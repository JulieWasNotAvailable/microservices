package handlers

import "github.com/gofiber/fiber/v2"

func Welcome(c *fiber.Ctx) error {
	// locals := c.Locals("user")
	return c.JSON(fiber.Map{
		"message": "welcome to user service!",
		"your context": c.Locals("oauth_token")})
}