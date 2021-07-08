package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/mohitsinghs/wormholes/auth"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/mohitsinghs/wormholes/links"
)

func Setup(linkStore links.Store, authStore auth.Store, factory *factory.Factory) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	store := session.New(session.Config{
		Expiration:   7 * 24 * time.Hour,
		KeyGenerator: utils.UUID,
	})

	app.Use(recover.New())

	authHandler := auth.NewHandler(store, authStore)
	linkHandler := links.NewHandler(store, linkStore, factory)

	app.Get("/:id", linkHandler.Redirect)

	apiV1 := app.Group("/api/v1")
	linkApi := apiV1.Group("/links")
	authApi := apiV1.Group("/auth")

	linkApi.Get("/:id", linkHandler.Get)
	linkApi.Put("/", linkHandler.Create)
	linkApi.Post("/:id", linkHandler.Update)
	linkApi.Delete("/:id", linkHandler.Delete)

	authApi.Post("/register", authHandler.Create)
	authApi.Get("/login", authHandler.Authenticate)
	authApi.Get("/logout", authHandler.Unauthenticate)
	authApi.Get("/user", authHandler.Get)
	authApi.Delete("/user", authHandler.Delete)

	return app
}
