package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/mohitsinghs/wormholes/links"
)

func Setup(linkStore links.Store, factory *factory.Factory) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	store := session.New(session.Config{
		Expiration:   7 * 24 * time.Hour,
		CookieName:   "worm",
		KeyGenerator: utils.UUID,
	})

	app.Use(recover.New())

	linkHandler := links.NewHandler(store, linkStore, factory)

	app.Get("/:id", linkHandler.Redirect)

	apiV1 := app.Group("/api/v1")
	linkApi := apiV1.Group("/links")

	linkApi.Get("/:id", linkHandler.Get)
	linkApi.Put("/", linkHandler.Create)
	linkApi.Post("/:id", linkHandler.Update)
	linkApi.Delete("/:id", linkHandler.Delete)

	return app
}
