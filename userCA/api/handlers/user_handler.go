package handlers

import (
	"errors"
	"strings"

	"net/http"

	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/bmmetadata"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/entities"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Username struct {
	Username string
}


// AddUser godoc
// @Summary Create a new user
// @Description Create a new user with default role (1)
// @Tags users
// @Accept json
// @Produce json
// @Param user body entities.User true "User data"
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 400 {object} presenters.UserErrorResponse
// @Failure 422 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /user [post]
func AddUser(service user.Service) fiber.Handler {
	return func (c *fiber.Ctx) error {
		var requestBody entities.User
		err := c.BodyParser(&requestBody)

		requestBody.RoleID = 1

		if requestBody.Email == "" || requestBody.Password == "" {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(errors.New(
				"email and password should not be null",
				)))
		}

		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateUserErrorResponse(err))
		}
		
		result, err := service.InsertUser(&requestBody)
		
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}
		return c.JSON(presenters.CreateUserSuccessResponse(result))
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /users [get]
func GetUsers(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.FetchUsers()

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		return c.JSON(presenters.CreateUsersSuccessResponse(result))
	}
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieve a single user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 400 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /userById/{id} [get]
func GetUserById(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuid, err := uuid.Parse(id)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		result, err := service.FetchUserById(uuid)
		if err != nil {
			
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		return c.JSON(presenters.CreateUserSuccessResponse2(result))
	}
}

// GetUserByEmail godoc
// @Summary Get user by email.
// @Description Retrieve a single user by their email. Requires authorization. Requires admin role.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email query string true "User Email"
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 400 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /userByEmail/ [get]
func GetUserByEmail(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Query("email")
		
		if email == "" {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateUserErrorResponse(errors.New(
				"please, specify email",
				)))
		}

		result, err := service.FetchUserByEmail(email)

		if err != nil{
			
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		return c.JSON(presenters.CreateUserSuccessResponse2(result))
	}
}

// GetBeatmakerByJWT godoc
// @Summary Get current user by JWT
// @Description Retrieve current user details from JWT token
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 401 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /user/me [get]
func GetUserByJWT(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.GetReqHeaders()
		authHeader := auth["Authorization"]
		splitToken := strings.Split(authHeader[0], "Bearer ")
		tokenStr := splitToken[1]

		token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(errors.New("couldn't parse")))
		}

		id := claims["id"].(string)
		uuid, err := uuid.Parse(id)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		user, err := service.FetchUserById(uuid)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		return c.JSON(presenters.CreateUserSuccessResponse2(user))
	}
}

// UpdateUser godoc
// @Summary Update user with metadata. Beatmaker role required. If user has no metadata, new is created with update data.
// @Description Update user details (requires authentication). Updates profile by jwt.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body presenters.User true "Updated user data"
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 400 {object} presenters.UserErrorResponse
// @Failure 401 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /user/me/withmeta [patch]
func UpdateBeatmaker(uservice user.Service, bmservice bmmetadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody presenters.User
		err := c.BodyParser(&requestBody)
		if err != nil {
			
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		uuid, err := uuid.Parse(id)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		user, err := uservice.FetchUserById(uuid)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}
		requestBody.RoleID = user.RoleID

		//can't save user data, if there is metadata struct inside it, so parse in two
		metadataBody := requestBody.Metadata
		emptyMetadata := presenters.Metadata{}
		requestBody.Metadata = emptyMetadata

		res, err := bmservice.FetchMetadataByUserId(uuid)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateUserErrorResponse(err))
		}

		if res == nil{
			emptyMetadataEntity := entities.Metadata{
				UserID : uuid,
			}
			_, err := bmservice.InsertMetadata(&emptyMetadataEntity)
			if err != nil{
				return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
			}
		}

		result, err := uservice.UpdateBeatmaker(uuid, &requestBody, &metadataBody)
		if err != nil {
			
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}
	
		return c.JSON(presenters.CreateUserSuccessResponse2(result))
	}
}

// UserIsBeatmaker godoc
// @Summary Upgrade user to beatmaker
// @Description Change user role to beatmaker (role 2). Requires jwt. Updates the role of the user based on given jwt.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /users/me/upgrade [get]
func UserIsBeatmaker (service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		id := claims["id"].(string)

		uuid, err := uuid.Parse(id)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateUserErrorResponse(err))
		}

		data := presenters.User{
			ID : uuid,
			RoleID: 2,
		}

		_, err = service.UpdateUser(&data)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		return c.JSON(&fiber.Map{
			"status": true,
			"data":   "updated successfully",
			"err":    nil,
		})	
	}
}

// UpdateUser godoc
// @Summary Update user WITHOUGHT METADATA.
// @Description Update user details (requires authentication). Updates profile by jwt.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body presenters.User true "Updated user data"
// @Success 200 {object} presenters.UserSuccessResponse
// @Failure 400 {object} presenters.UserErrorResponse
// @Failure 401 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /user/me [patch]
func UpdateUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := presenters.User{}
		err := c.BodyParser(&requestBody)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		uuid, err := uuid.Parse(id)
		if err !=nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		requestBody.ID = uuid
		result, err := service.UpdateUser(&requestBody)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateUserSuccessResponse2(result))
	}
}

// RemoveUser godoc
// @Summary Delete user
// @Description Delete the current user account by jwt.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} presenters.UserErrorResponse
// @Failure 500 {object} presenters.UserErrorResponse
// @Router /users/me [delete]
func RemoveUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		id := claims["id"].(string)

		uuid, err := uuid.Parse(id)
		
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		err = service.RemoveUser(uuid)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateUserErrorResponse(err))
		}

		return c.JSON(&fiber.Map{
			"status": true,
			"data":   "updated successfully",
			"err":    nil,
		})		
	}
}

func PostBeatMock(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{
		"status": true,
		"data":   "congrats! you posted the beat!",
		"err":    nil,
	})
}