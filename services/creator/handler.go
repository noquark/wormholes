package creator

import (
	"context"
	"reflect"
	"wormholes/services/creator/reserve"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Fiber route handlers for link.
type Handler struct {
	backend  Store
	ingestor *Ingestor
	cache    *redis.Client
	reserve  reserve.Reserve
}

func NewHandler(
	backend Store,
	ingestor *Ingestor,
	cache *redis.Client,
	reserve reserve.Reserve,
) *Handler {
	return &Handler{
		backend,
		ingestor,
		cache,
		reserve,
	}
}

func (h *Handler) Setup(linkAPI fiber.Router) {
	linkAPI.Get("/:id", h.Get)
	linkAPI.Put("/", h.Create)
	linkAPI.Post("/:id", h.Update)
	linkAPI.Delete("/:id", h.Delete)
}

type LinkCreateRequest struct {
	Tag    string `json:"tag"`
	Target string `json:"target"`
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req LinkCreateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("create: failed to parsing request")

		return fiber.ErrBadRequest
	}

	var link *Link

	id, err := h.reserve.GetID()
	if err != nil {
		log.Error().Err(err).Msg("create: failed to get id")

		return fiber.ErrInternalServerError
	}

	link = NewLink(id, req.Target, req.Tag)
	h.ingestor.Push(link)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Link Created",
		"id":     link.ID,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var link Link
	if err := c.BodyParser(&link); err != nil {
		log.Error().Err(err).Msg("error parsing request")

		return fiber.ErrBadRequest
	}

	if err := h.backend.Update(&link); err != nil {
		log.Error().Err(err).Msg("error updating link")

		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return fiber.ErrBadRequest
	}

	var link Link

	err := h.cache.HGetAll(context.Background(), id).Scan(&link)
	if err != nil || reflect.ValueOf(link).IsZero() {
		log.Err(err).Msg("get: cache miss")

		// If key does not exists, query db
		link, err = h.backend.Get(id)

		if err != nil {
			log.Error().Err(err).Msg("get: error getting link")

			return fiber.ErrBadRequest
		}

		err = h.cache.HSet(context.Background(), id, "id", link.ID, "target", link.Target, "tag", link.Tag).Err()
		if err != nil {
			log.Warn().Err(err).Msg("get: failed to cache")
		}

		return c.Status(fiber.StatusOK).JSON(link)
	}

	return c.Status(fiber.StatusOK).JSON(link)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return fiber.ErrBadRequest
	}

	if err := h.backend.Delete(id); err != nil {
		log.Error().Err(err).Msg("error deleting link")

		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}
