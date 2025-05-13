package presenters

import (
	"github.com/JulieWasNotAvailable/microservices/user/internal/entities"
	"github.com/google/uuid"
)

// Metadata represents metadata information
// @Description Metadata entity containing social links and description
type Metadata struct {
    ID               uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
    VkUrl           string    `json:"vkUrl,omitempty" example:"https://vk.com/username"`
    TelegramUrl     string    `json:"telegramUrl,omitempty" example:"https://t.me/username"`
    InstagramUrl    string    `json:"instagramUrl,omitempty" example:"https://instagram.com/username"`
    Description     string    `json:"description,omitempty" example:"Artist profile description"`
    SubscriptionTypeID int    `json:"subscriptionTypeId,omitempty" example:"1"`
    UserID          uuid.UUID `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// MetadataResponse represents a successful metadata response
// @Description Standard response for single metadata operation
type MetadataSuccessResponse struct {
    Status bool     `json:"status" example:"true"`
    Data   Metadata `json:"data"`
    Error  *string  `json:"error" example:"null"`
}

// MetadataListResponse represents a successful metadata list response
// @Description Standard response for multiple metadata items
type MetadataListResponse struct {
    Status bool       `json:"status" example:"true"`
    Data   []Metadata `json:"data"`
    Error  *string    `json:"error" example:"null"`
}

// MetadataErrorResponse represents an error response
// @Description Error response for metadata operations
type MetadataErrorResponse struct {
    Status bool    `json:"status" example:"false"`
    Data   *string `json:"data" example:"null"`
    Error  string  `json:"error" example:"error message"`
}

// CreateMetadataSuccessResponse creates a success response from entities.Metadata
// @Summary Creates success response
// @Description Converts entities.Metadata to API success response
func CreateMetadataSuccessResponse(data *entities.Metadata) *MetadataSuccessResponse {
    metadata := Metadata{
        ID: data.ID,
        VkUrl:            data.VkUrl,
        TelegramUrl:      data.TelegramUrl,
        InstagramUrl:     data.InstagramUrl,
        Description:      data.Description,
        SubscriptionTypeID: data.SubscriptionTypeID,
        UserID:          data.UserID,
    }

    return &MetadataSuccessResponse{
        Status: true,
        Data:   metadata,
        Error:  nil,
    }
}

// CreateMetadataResponse creates a success response from Metadata
// @Summary Creates metadata response
// @Description Converts Metadata to API success response
func CreateMetadataSuccessResponse2(data *Metadata) *MetadataSuccessResponse {
    return &MetadataSuccessResponse{
        Status: true,
        Data:   *data,
        Error:  nil,
    }
}

// CreateMetadataListResponse creates a success response for metadata list
// @Summary Creates metadata list response
// @Description Converts []Metadata to API success response
func CreateMetadataListResponse(data *[]Metadata) *MetadataListResponse {
    return &MetadataListResponse{
        Status: true,
        Data:   *data,
        Error:  nil,
    }
}

// CreateMetadataErrorResponse creates an error response
// @Summary Creates error response
// @Description Converts error to standard error response
func CreateMetadataErrorResponse(err error) *MetadataErrorResponse {
    return &MetadataErrorResponse{
        Status: false,
        Data:   nil,
        Error:  err.Error(),
    }
}