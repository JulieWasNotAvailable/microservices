package routes

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/guards"
	"github.com/JulieWasNotAvailable/microservices/user/internal/bmmetadata"
	"github.com/JulieWasNotAvailable/microservices/user/internal/user"
	"github.com/JulieWasNotAvailable/microservices/user/internal/authentication"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service, bmservice bmmetadata.Service) {
	app.Post("/user", handlers.AddUser(service))
	app.Get("/users", handlers.GetUsers(service))
	app.Get("/userById/:id", handlers.GetUserById(service))
	app.Get("/userByEmail/", guards.ProtectedRequiresModerator(), handlers.GetUserByEmail(service))
	app.Delete("/users/me", guards.Protected(), handlers.RemoveUser(service))

	app.Get("/users/me/upgrade", guards.Protected(), handlers.UserIsBeatmaker(service))
	app.Patch("/user/me", guards.Protected(), handlers.UpdateUser(service))
	app.Patch("/user/me/withmeta", guards.ProtectedRequiresBeatmaker(), handlers.UpdateBeatmaker(service, bmservice))
	app.Get("/user/me", guards.Protected(), handlers.GetUserByJWT(service))
	
	app.Post("/login", authentication.Login(service))
	app.Post("/register", authentication.Register(service))
	app.Post("/postNewBeatMock", guards.ProtectedRequiresBeatmaker(), handlers.PostBeatMock)
}
