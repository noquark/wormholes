package links

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mohitsinghs/wormholes/constants"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/mohitsinghs/wormholes/pipe"
	"github.com/rs/zerolog/log"
)

// Fiber route handlers for link

type Handler struct {
	backend  Store
	factory  *factory.Factory
	pipe     *pipe.Pipe
	ingestor *Ingestor
}

func NewHandler(backend Store, factory *factory.Factory, pipe *pipe.Pipe, ingestor *Ingestor) *Handler {
	return &Handler{
		backend,
		factory,
		pipe,
		ingestor,
	}
}

type LinkCreateRequest struct {
	Custom string `json:"custom"`
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
	if req.Custom != "" {
		if h.factory.ID.Exists([]byte(req.Custom)) {
			log.Error().Msg("link already exists")
			return fiber.ErrBadRequest
		}
		link = New(req.Custom, req.Target, req.Tag)
		h.factory.ID.Add([]byte(req.Custom))
	} else {
		id, err := h.factory.NewId()
		if err != nil {
			log.Error().Err(err).Msg("error generating id")
			return fiber.ErrInternalServerError
		}
		link = New(id, req.Target, req.Tag)
	}

	h.ingestor.Push(link)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Link Created",
		"id":     link.Id,
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
	link, err := h.backend.Get(id)
	if err != nil {
		log.Error().Err(err).Msg("error getting link")
		return fiber.ErrBadRequest
	}
	return c.Status(fiber.StatusOK).JSON(link)
}

func (h *Handler) Redirect(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return fiber.ErrBadRequest
	}
	if !h.factory.ID.Exists([]byte(id)) {
		return fiber.ErrNotFound
	}
	link, err := h.backend.Get(id)
	if err != nil {
		log.Error().Err(err).Msg("error getting link")
		return fiber.ErrInternalServerError
	}
	var cookie string
	if c.Cookies(constants.COOKIE_NAME) == "" {
		cookie := h.factory.NewCookie()
		c.Cookie(&fiber.Cookie{
			Name:    constants.COOKIE_NAME,
			Value:   cookie,
			Expires: time.Now().Add(time.Hour * 24 * 180),
		})
	} else {
		cookie = c.Cookies(constants.COOKIE_NAME)
	}
	c.Set(fiber.HeaderCacheControl, constants.CACHE_CONTROL)
	h.pipe.Push(pipe.NewEvent(link.Id, link.Tag, cookie, c.Get(fiber.HeaderUserAgent), c.IP()))
	return c.Redirect(link.Target, fiber.StatusMovedPermanently)
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
