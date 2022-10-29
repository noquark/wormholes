package creator

import (
	"context"
	"reflect"
	"wormholes/services/creator/ingestor"
	"wormholes/services/creator/links"
	"wormholes/services/creator/reserve"
	"wormholes/services/creator/store"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Fiber route handlers for link.
type Handler struct {
	backend  store.Store
	ingestor *ingestor.Ingestor
	cache    *redis.Client
	reserve  reserve.Reserve
}

func NewHandler(
	backend store.Store,
	ingestor *ingestor.Ingestor,
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

func (h *Handler) Create(ctx *fiber.Ctx) error {
	var req LinkCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("create: failed to parsing request")

		return fiber.ErrBadRequest
	}

	var link *links.Link

	newID, err := h.reserve.GetID()
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

	err := h.cache.HGetAll(context.Background(), shortID).Scan(&link)
	if err != nil || reflect.ValueOf(link).IsZero() {
		log.Err(err).Msg("get: cache miss")

		// If key does not exists, query db
		link, err = h.backend.Get(shortID)

		if err != nil {
			log.Error().Err(err).Msg("get: error getting link")

			return fiber.ErrBadRequest
		}

		err = h.cache.HSet(context.Background(), shortID, "id", link.ID, "target", link.Target, "tag", link.Tag).Err()
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
