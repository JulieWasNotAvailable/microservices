package presenters

// import "github.com/gofiber/fiber/v2"

// @Description Data Token
type Data struct {
	Token string `json:"token"`
}

// type GoogleResponse struct {
// 	ID       string `json:"id"`
// 	Email    string `json:"email"`
// 	Verified bool   `json:"verified_email"`
// 	Picture  string `json:"picture"`
// }

// func GoogleOauthErrorResponse(err error) *fiber.Map {
// 	return &fiber.Map{
// 		"status": false,
// 		"data":   "",
// 		"error":  err.Error(),
// 	}
// }