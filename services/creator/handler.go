package main

import (
	"context"
	"fmt"
	"wormholes/protos"

	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
)

// Fiber route handlers for link.
type Handler struct {
	backend      Store
	ingestor     *Ingestor
	cache        *redis.Pool
	bucket       *protos.Bucket
	bucketClient protos.BucketServiceClient
}

func NewHandler(
	backend Store,
	ingestor *Ingestor,
	cache *redis.Pool,
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

	id := h.bucket.Ids[0]
	h.bucket.Ids = h.bucket.Ids[1:]

	return id, nil
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

	var link *Link

	// retrieve data from cache first
	conn := h.cache.Get()
	val, err := redis.Values(conn.Do("HGETALL", id))
	errScan := redis.ScanStruct(val, link)

	if err != nil || errScan != nil {
		log.Err(err).Msg("get: cache miss")

		// If key does not exists, query db
		link, err := h.backend.Get(id)
		if err != nil {
			log.Error().Err(err).Msg("get: error getting link")

			return fiber.ErrBadRequest
		}

		_, _ = conn.Do("HSET", redis.Args{}.Add(id).AddFlat(link)...)
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
