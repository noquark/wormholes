package main

import (
	"context"
	"errors"
	"fmt"
	"wormholes/protos"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var ErrNoIds = errors.New("failed to get id")

// Fiber route handlers for link.
type Handler struct {
	backend      Store
	ingestor     *Ingestor
	cache        *redis.Client
	bucket       *protos.Bucket
	bucketClient protos.BucketServiceClient
}

func NewHandler(
	backend Store,
	ingestor *Ingestor,
	cache *redis.Client,
	bucketClient protos.BucketServiceClient,
) *Handler {
	return &Handler{
		backend,
		ingestor,
		cache,
		&protos.Bucket{},
		bucketClient,
	}
}

func (h *Handler) fetchBucket() error {
	bucket, err := h.bucketClient.GetBucket(context.Background(), &protos.Empty{})
	if err != nil {
		return fmt.Errorf("grpc: failed to fetch bucket: %w", err)
	}

	h.bucket = bucket

	return nil
}

func (h *Handler) getID() (string, error) {
	if len(h.bucket.Ids) == 0 {
		err := h.fetchBucket()
		if err != nil {
			log.Warn().Err(err)

			return "", err
		}
	}

	if len(h.bucket.Ids) > 0 {
		id := h.bucket.Ids[0]
		h.bucket.Ids = h.bucket.Ids[1:]

		return id, nil
	}

	return "", ErrNoIds
}

func (h *Handler) Setup(app *fiber.App) {
	// group routes
	apiV1 := app.Group("/api/v1")
	linkAPI := apiV1.Group("/links")

	// link routes
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
		log.Error().Err(err).Msg("error parsing request")

		return fiber.ErrBadRequest
	}

	var link *Link

	id, err := h.getID()
	if err != nil {
		return err
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

	var link *Link

	ctx := context.Background()

	err := h.cache.HGetAll(ctx, id).Scan(link)
	if err != nil {
		log.Err(err).Msg("get: cache miss")

		// If key does not exists, query db
		link, err := h.backend.Get(id)
		if err != nil {
			log.Error().Err(err).Msg("get: error getting link")

			return fiber.ErrBadRequest
		}

		_ = h.cache.HSet(ctx, id, link).Err()
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
