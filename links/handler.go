package links

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mohitsinghs/wormholes/constants"
	"github.com/mohitsinghs/wormholes/pipe"
	"github.com/mohitsinghs/wormholes/state"
)

// Fiber route handlers for link

type Handler struct {
	backend Store
	state   *state.State
}

func NewHandler(backend Store, state *state.State) *Handler {
	return &Handler{
		backend,
		state,
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
		if h.state.Factory.Exists(req.Custom) {
			log.Println("Link already exists")
			return fiber.ErrBadRequest
		}
		link = New(req.Custom, req.Target, req.Tag)
		h.state.Factory.Add(req.Custom)
	} else {
		id, err := h.state.Factory.NewId()
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
	if !h.state.Factory.Exists(id) {
		return fiber.ErrNotFound
	}
	link, err := h.backend.Get(id)
	if err != nil {
		log.Printf("Error getting link : %v", err)
		return fiber.ErrInternalServerError
	}
	var cookie string
	if c.Cookies(constants.COOKIE_NAME) == "" {
		cookie := h.state.Factory.NewCookie()
		c.Cookie(&fiber.Cookie{
			Name:    constants.COOKIE_NAME,
			Value:   cookie,
			Expires: time.Now().Add(time.Hour * 24 * 180),
		})
	} else {
		cookie = c.Cookies(constants.COOKIE_NAME)
	}
	c.Set(fiber.HeaderCacheControl, constants.CACHE_CONTROL)
	h.state.Pipe.Push(pipe.Event{
		Time:   time.Now(),
		Link:   link.Id,
		Tag:    link.Tag,
		Cookie: cookie,
		UA:     c.Get(fiber.HeaderUserAgent),
		IP:     c.IP(),
	})
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
