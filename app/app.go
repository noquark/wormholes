package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/mohitsinghs/wormholes/auth"
	"github.com/mohitsinghs/wormholes/links"
	"github.com/mohitsinghs/wormholes/state"
)

func Setup(state *state.State) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
	})

	store := session.New(session.Config{
		Expiration:   7 * 24 * time.Hour,
		KeyGenerator: state.Factory.NewCookie,
	})

	authStore := auth.NewStore(state.DB)
	linkStore := links.NewStore(state.DB)

	authHandler := auth.NewHandler(store, authStore)
	linkHandler := links.NewHandler(linkStore, state)

	state.Factory.TryRestore(linkStore.Ids)
	authHandler.EnsureDefault(
		state.Config.Admin.Email,
		state.Config.Admin.Password,
	)

	app.Use(etag.New())
	app.Use(recover.New())

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
