package links

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/mohitsinghs/wormholes/internal/factory"
)

// Fiber route handlers for link

type Handler struct {
	store   *session.Store
	backend Store
	factory *factory.Factory
}

func NewHandler(store *session.Store, backend Store, factory *factory.Factory) *Handler {
	return &Handler{
		store,
		backend,
		factory,
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
		log.Printf("Error parsing req : %v", err)
		return fiber.ErrBadRequest
	}

	var link *Link
	if req.Custom != "" {
		link = New(req.Custom, req.Target, req.Tag)
		h.factory.Add(req.Custom)
	} else {
		id, err := h.factory.NewId()
		if err != nil {
			log.Printf("Error generating id : %v", err)
			return fiber.ErrInternalServerError
		}
		link = New(id, req.Target, req.Tag)
	}

	if err := h.backend.Insert(link); err != nil {
		log.Printf("Error inserting link : %v", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Link Created",
		"id":     link.Id,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var link Link
	if err := c.BodyParser(&link); err != nil {
		log.Printf("Error parsing req : %v", err)
		return fiber.ErrBadRequest
	}
	if err := h.backend.Update(&link); err != nil {
		log.Printf("Error updating link : %v", err)
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
		log.Printf("Error getting link : %v", err)
		return fiber.ErrBadRequest
	}
	return c.Status(fiber.StatusOK).JSON(link)
}

func (h *Handler) Redirect(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return fiber.ErrBadRequest
	}
	if !h.factory.Exists(id) {
		return fiber.ErrNotFound
	}
	link, err := h.backend.Get(id)
	if err != nil {
		log.Printf("Error getting link : %v", err)
		return fiber.ErrInternalServerError
	}
	return c.Redirect(link.Target, fiber.StatusMovedPermanently)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return fiber.ErrBadRequest
	}

	if err := h.backend.Delete(id); err != nil {
		log.Printf("Error deleting link : %v", err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}
