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

// nolint:funlen
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
	linkAPI := apiV1.Group("/links")
	authAPI := apiV1.Group("/auth")
	statsAPI := apiV1.Group("/stats")

	// configure middlewares on routes
	csrfMW := csrf.New(csrf.Config{
		KeyGenerator: factory.NewCookie,
	})
	authAPI.Use(csrfMW)
	statsAPI.Use(csrfMW)
	statsAPI.Use(authHandler.VerifyAdmin)

	// link routes
	linkAPI.Get("/:id", linkHandler.Get)
	linkAPI.Put("/", linkHandler.Create)
	linkAPI.Post("/:id", linkHandler.Update)
	linkAPI.Delete("/:id", linkHandler.Delete)

	// auth routes
	authAPI.Post("/register", authHandler.Create)
	authAPI.Get("/login", authHandler.Authenticate)
	authAPI.Get("/logout", authHandler.Unauthenticate)
	authAPI.Get("/user", authHandler.Get)
	authAPI.Delete("/user", authHandler.Delete)

	// stats routes
	statsAPI.Get("/cards", statsHandler.Cards)

	return app
}
