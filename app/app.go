package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/mohitsinghs/wormholes/auth"
	"github.com/mohitsinghs/wormholes/config"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/mohitsinghs/wormholes/links"
	"github.com/mohitsinghs/wormholes/pipe"
	"github.com/mohitsinghs/wormholes/stats"
)

func Setup(config *config.Config, factory *factory.Factory, pipe *pipe.Pipe) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
	})

	store := session.New(session.Config{
		Expiration:   7 * 24 * time.Hour,
		KeyGenerator: factory.NewCookie,
	})

	db := config.Postgres.Connect()

	authStore := auth.NewStore(db)
	linkStore := links.NewStore(db)
	statsStore := stats.NewStore(db)

	authHandler := auth.NewHandler(store, authStore)
	linkHandler := links.NewHandler(linkStore, factory, pipe)
	statsHandler := stats.NewHandler(statsStore)

	factory.TryRestore(linkStore.Ids)
	authHandler.EnsureDefault(config.Admin)

	app.Use(etag.New())
	app.Use(recover.New())

	app.Get("/:id", linkHandler.Redirect)

	apiV1 := app.Group("/api/v1")
	linkApi := apiV1.Group("/links")
	authApi := apiV1.Group("/auth")
	statsApi := apiV1.Group("/stats")

	statsApi.Use(authHandler.VerifyAdmin)

	linkApi.Get("/:id", linkHandler.Get)
	linkApi.Put("/", linkHandler.Create)
	linkApi.Post("/:id", linkHandler.Update)
	linkApi.Delete("/:id", linkHandler.Delete)

	authApi.Post("/register", authHandler.Create)
	authApi.Get("/login", authHandler.Authenticate)
	authApi.Get("/logout", authHandler.Unauthenticate)
	authApi.Get("/user", authHandler.Get)
	authApi.Delete("/user", authHandler.Delete)

	statsApi.Get("/cards", statsHandler.Cards)

	return app
}
