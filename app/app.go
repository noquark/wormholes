package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
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

func Setup(config *config.Config, factory *factory.Factory) *fiber.App {
	//  shared connection pool
	db := config.Postgres.Connect()
	tsdb := config.Timescale.Connect()
	p := pipe.New(config, tsdb).Start().Wait()
	i := links.NewIngestor(db, config.Pipe.BatchSize).Start()

	// session store
	store := session.New(session.Config{
		Expiration:   7 * 24 * time.Hour,
		KeyGenerator: factory.NewCookie,
	})

	// db stores for routes
	authStore := auth.NewStore(db)
	linkStore := links.NewStore(db)
	statsStore := stats.NewStore(db, tsdb)

	// route handlers
	authHandler := auth.NewHandler(store, authStore)
	linkHandler := links.NewHandler(linkStore, factory, p, i)
	statsHandler := stats.NewHandler(statsStore)

	// restore bloom-filters
	factory.Restore(linkStore.Ids)

	// ensure existence of admin
	authHandler.EnsureDefault(config.Admin)

	// create app
	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
	})

	// configure middlewares globally
	app.Use(etag.New())
	app.Use(recover.New())

	// handle global route
	app.Get("/:id", linkHandler.Redirect)

	// group routes
	apiV1 := app.Group("/api/v1")
	linkApi := apiV1.Group("/links")
	authApi := apiV1.Group("/auth")
	statsApi := apiV1.Group("/stats")

	// configure middlewares on routes
	csrfMW := csrf.New(csrf.Config{
		KeyGenerator: factory.NewCookie,
	})
	authApi.Use(csrfMW)
	statsApi.Use(csrfMW)
	statsApi.Use(authHandler.VerifyAdmin)

	// link routes
	linkApi.Get("/:id", linkHandler.Get)
	linkApi.Put("/", linkHandler.Create)
	linkApi.Post("/:id", linkHandler.Update)
	linkApi.Delete("/:id", linkHandler.Delete)

	// auth routes
	authApi.Post("/register", authHandler.Create)
	authApi.Get("/login", authHandler.Authenticate)
	authApi.Get("/logout", authHandler.Unauthenticate)
	authApi.Get("/user", authHandler.Get)
	authApi.Delete("/user", authHandler.Delete)

	// stats routes
	statsApi.Get("/cards", statsHandler.Cards)

	return app
}
