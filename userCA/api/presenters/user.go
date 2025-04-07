package presenters

import (
	"github.com/JulieWasNotAvailable/microservices/user/pkg/entities"
	"github.com/google/uuid"
)

// User represents the user response structure
// @Description User information
type User struct {
	ID              uuid.UUID    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email           string `json:"email,omitempty" example:"user@example.com"`
	Password        string `json:"-"` // Hidden from JSON responses
	Username        string `json:"username,omitempty" example:"johndoe"`
	Firstname       string `json:"firstname,omitempty" example:"John"`
	Lastname        string `json:"lastname,omitempty" example:"Doe"`
	Patronymic      string `json:"patronymic,omitempty" example:"Smith"`

	RoleID          uint   `json:"roleId,omitempty" example:"1"`
	SubscriptionID  int   `json:"subscriptionId,omitempty" example:"1"`
	
	UsersFavourites int `json:"usersFavourites,omitempty" example:"1"`
	FollowerOf      int `json:"followerOf,omitempty" example:"4"`
	Metadata  	 	Metadata   `gorm:"foreignKey:UserID"`
}

// UserSuccessResponse represents a successful user response
// @Description Success response containing user data
type UserSuccessResponse struct {
	Status bool `json:"status" example:"true"`
	Data   User `json:"data"`
	Error  *string `json:"error" example:"null"`
}

// UserListSuccessResponse represents a successful list user response
// @Description Success response containing multiple user items
type UserListSuccessResponse struct {
	Status bool   `json:"status" example:"true"`
	Data   []User `json:"data"`
	Error  *string `json:"error" example:"null"`
}

// UserErrorResponse represents an error response
// @Description Error response for user operations
type UserErrorResponse struct {
	Status bool    `json:"status" example:"false"`
	Data   *string `json:"data" example:"null"`
	Error  string  `json:"error" example:"error message"`
}

func CreateUserSuccessResponse(data *entities.User) *UserSuccessResponse {
	user := User{
		ID:              data.ID,
		Username:        data.Username,
		Password:        data.Password, // Never expose password
		RoleID:          data.RoleID,
		SubscriptionID:  data.SubscriptionID,
		Email:           data.Email,
		UsersFavourites: data.UsersFavourites,
		FollowerOf:      data.FollowerOf,
	}

	return &UserSuccessResponse{
		Status: true,
		Data:   user,
		Error:  nil,
	}
}

func CreateUserSuccessResponse2(data *User) *UserSuccessResponse {
	return &UserSuccessResponse{
		Status: true,
		Data:   *data,
		Error:  nil,
	}
}

func CreateUsersSuccessResponse(data *[]User) *UserListSuccessResponse {
	return &UserListSuccessResponse{
		Status: true,
		Data:   *data,
		Error:  nil,
	}
}

func CreateUserErrorResponse(err error) *UserErrorResponse {
	return &UserErrorResponse{
		Status: false,
		Data:   nil,
		Error:  err.Error(),
	}
}