package main

import (
	_ "embed"
	"reflect"
	"time"
	"wormholes/ingestor"
	"wormholes/internal/cache"
	"wormholes/internal/links"
	"wormholes/ipc"
	"wormholes/store"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/noquark/nanoid"
	"github.com/rs/zerolog/log"
)

// Fiber route handlers for link.
type Handler struct {
	backend  store.Store
	ingestor *ingestor.Ingestor
	cache    *cache.Cache
	store    *ipc.Store
}

const (
	CookieExpiryTime = time.Hour * 24 * 180
	CacheControl     = "private, max-age=90"
	CookieName       = "_wh"
	MaxTry           = 10
	CookieSize       = 21
	backOffTime      = 5e3
)

func NewHandler(
	backend store.Store,
	in *ingestor.Ingestor,
	cache *cache.Cache,
	ipcStore *ipc.Store,
) *Handler {
	return &Handler{
		backend,
		in,
		cache,
		ipcStore,
	}
}

// Generate a random cookie with retry on failure.
func NewCookie() string {
	cookie, err := nanoid.New(CookieSize)
	if err != nil {
		for i := 0; i < MaxTry; i++ {
			cookie, err = nanoid.New(CookieSize)
			if err != nil {
				continue
			}

			break
		}
	}

	return cookie
}

func (h *Handler) Setup(app fiber.Router) {
	app.Get("/:id", h.Redirect)

	api := app.Group("api")
	api.Get("/:id", h.Get)
	api.Put("/", h.Create)
	api.Post("/:id", h.Update)
	api.Delete("/:id", h.Delete)
}

type LinkCreateRequest struct {
	Tag    string `json:"tag"`
	Target string `json:"target"`
}

func (h *Handler) Create(ctx *fiber.Ctx) error {
	var req LinkCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("create: failed to parsing request")

		return fiber.ErrBadRequest
	}

	var link *links.Link

	newID, err := h.store.GetID()
	if err != nil {
		log.Error().Err(err).Msg("create: failed to get id")

		return fiber.ErrInternalServerError
	}

	link = links.New(newID, req.Target, req.Tag)
	h.ingestor.Push(link)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Link Created",
		"id":     link.ID,
	})
}

func (h *Handler) Update(ctx *fiber.Ctx) error {
	var link links.Link
	if err := ctx.BodyParser(&link); err != nil {
		log.Error().Err(err).Msg("error parsing request")

		return fiber.ErrBadRequest
	}

	if err := h.backend.Update(&link); err != nil {
		log.Error().Err(err).Msg("error updating link")

		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) Get(ctx *fiber.Ctx) error {
	shortID := ctx.Params("id")
	if len(shortID) == 0 {
		return fiber.ErrBadRequest
	}

	var link links.Link

	err := h.cache.GetLink(&link, shortID)
	if err != nil || reflect.ValueOf(link).IsZero() {
		log.Err(err).Msg("get: cache miss")

		// If key does not exists, query db
		link, err = h.backend.Get(shortID)

		if err != nil {
			log.Error().Err(err).Msg("get: error getting link")

			return fiber.ErrBadRequest
		}

		err = h.cache.SetLink(link, shortID)
		if err != nil {
			log.Warn().Err(err).Msg("get: failed to cache")
		}

		return ctx.Status(fiber.StatusOK).JSON(link)
	}

	return ctx.Status(fiber.StatusOK).JSON(link)
}

func (h *Handler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) == 0 {
		return fiber.ErrBadRequest
	}

	if err := h.backend.Delete(id); err != nil {
		log.Error().Err(err).Msg("error deleting link")

		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) Redirect(c *fiber.Ctx) error {
	shortID := c.Params("id")
	if len(shortID) == 0 {
		return fiber.ErrBadRequest
	}

	var link links.Link

	err := h.cache.GetLink(&link, shortID)
	if err != nil || reflect.ValueOf(link).IsZero() {
		log.Err(err).Msg("redirect: cache miss")

		// If key does not exists, query db
		link, err := h.backend.Get(shortID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return fiber.ErrNotFound
			}
			log.Error().Err(err).Msg("redirect: error getting link")

			return fiber.ErrInternalServerError
		}

		err = h.cache.SetLink(link, shortID)
		if err != nil {
			log.Warn().Err(err).Msg("redirect: failed to cache")
		}
	}

	if c.Cookies(CookieName) == "" {
		cookie := NewCookie()

		c.Cookie(&fiber.Cookie{
			Name:    CookieName,
			Value:   cookie,
			Expires: time.Now().Add(CookieExpiryTime),
		})
	}

	c.Set(fiber.HeaderCacheControl, CacheControl)

	return c.Redirect(link.Target, fiber.StatusMovedPermanently)
}
