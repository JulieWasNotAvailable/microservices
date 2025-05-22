package presenters

// @Description presenters.SuccessResponse
type SuccessResponse struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:""`
}

// @Description presenters.ListResponse
type ListResponse struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:""`
}

// @Description presenters.ErrorResponse
type ErrorResponse struct {
	Status bool   `json:"status" example:"false"`
	Data   string `json:"data" example:""`
	Error  string `json:"error" example:"error message"`
}

func CreateSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Status: true,
		Data:   data,
		Error:  "",
	}
}

func CreateListResponse(data interface{}) *ListResponse {
	return &ListResponse{
		Status: true,
		Data:   data,
		Error:  "",
	}
}

func CreateErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Status: false,
		Data:   "",
		Error:  err.Error(),
	}
}
