package stats

import "github.com/gofiber/fiber/v2"

type Handler struct {
	backend Store
}

func NewHandler(backend Store) *Handler {
	return &Handler{
		backend,
	}
}

func (h *Handler) Cards(c *fiber.Ctx) error {
	overview, err := h.backend.Overview()
	if err != nil {
		return fiber.ErrInternalServerError
	}
	dbSize, err := h.backend.DBSize()
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"overview": overview,
		"db_size":  dbSize,
	})
}
