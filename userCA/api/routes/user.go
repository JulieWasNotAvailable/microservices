package routes

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/user/internal/middleware"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/bmmetadata"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service, bmservice bmmetadata.Service) {
	app.Post("/user", handlers.AddUser(service))
	app.Get("/users", handlers.GetUsers(service))
	app.Get("/userById/:id", handlers.GetUserById(service))
	app.Get("/userByEmail/", middleware.ProtectedRequiresModerator(), handlers.GetUserByEmail(service))
	app.Delete("/users/me", middleware.Protected(), handlers.RemoveUser(service))

	app.Get("/users/me/upgrade", middleware.Protected(), handlers.UserIsBeatmaker(service))
	app.Patch("/user/me", middleware.Protected(), handlers.UpdateUser(service))
	app.Patch("/user/me/withmeta", middleware.ProtectedRequiresBeatmaker(), handlers.UpdateBeatmaker(service, bmservice))
	app.Get("/user/me", middleware.Protected(), handlers.GetUserByJWT(service))
	
	app.Post("/login", middleware.Login(service))
	app.Post("/register", middleware.Register(service))
	app.Post("/postNewBeatMock", middleware.ProtectedRequiresBeatmaker(), handlers.PostBeatMock)
}
