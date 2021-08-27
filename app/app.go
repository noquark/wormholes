package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mohitsinghs/wormholes/config"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/mohitsinghs/wormholes/links"
	"github.com/mohitsinghs/wormholes/pipe"
)

func Setup(config *config.Config, factory *factory.Factory) *fiber.App {
	//  shared connection pool
	db := config.Postgres.Connect()
	tsdb := config.Timescale.Connect()
	p := pipe.New(config, tsdb).Start().Wait()
	i := links.NewIngestor(db, config.Pipe.BatchSize).Start()

	// db stores for routes
	linkStore := links.NewStore(db)

	// route handlers
	linkHandler := links.NewHandler(linkStore, factory, p, i)

	// restore bloom-filters
	factory.Restore(linkStore.Ids)

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

	// link routes
	linkAPI.Get("/:id", linkHandler.Get)
	linkAPI.Put("/", linkHandler.Create)
	linkAPI.Post("/:id", linkHandler.Update)
	linkAPI.Delete("/:id", linkHandler.Delete)

	return app
}
