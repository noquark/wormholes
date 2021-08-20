package auth

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/mohitsinghs/wormholes/config"
)

type Handler struct {
	store   *session.Store
	backend Store
}

func NewHandler(store *session.Store, backend Store) *Handler {
	return &Handler{
		store,
		backend,
	}
}

type CreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

func (h *Handler) EnsureDefault(admin config.Admin) {
	hash, err := GenerateFromPassword([]byte(admin.Password))
	if err != nil {
		log.Panicln("failed to hash password for default user : ", err)
	}

	user := New(admin.Email, true)

	err = h.backend.InsertSafe(user, hash)
	if err != nil {
		log.Panicln(err)
	}
}

// Create a new user.
func (h *Handler) Create(c *fiber.Ctx) error {
	s, err := h.store.Get(c)
	if err != nil || !s.Get("admin").(bool) {
		return fiber.ErrUnauthorized
	}

	var createReq CreateRequest

	if err := c.BodyParser(&createReq); err != nil {
		log.Printf("Error parsing req : %v", err)

		return fiber.ErrBadRequest
	}

	hash, err := GenerateFromPassword([]byte(createReq.Password))
	if err != nil {
		log.Printf("Error hasing password : %v", err)

		return fiber.ErrBadRequest
	}

	user := New(createReq.Email, createReq.IsAdmin)
	if err = h.backend.Insert(user, hash); err != nil {
		log.Printf("Error creating user : %v", err)

		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// Authenticate user based on basic auth credentials.
func (h *Handler) Authenticate(c *fiber.Ctx) error {
	auth := c.Get(fiber.HeaderAuthorization)
	if auth == "" {
		return fiber.ErrUnauthorized
	}

	email, password, ok := ParseAuth(auth)
	if !ok {
		return fiber.ErrUnauthorized
	}

	isValid, isAdmin := h.backend.ValidateAuth(email, password)

	if isValid {
		session, err := h.store.Get(c)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		session.Set("email", email)
		session.Set("admin", isAdmin)

		err = session.Save()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.SendStatus(fiber.StatusOK)
	}

	return fiber.ErrUnauthorized
}

func (h *Handler) Unauthenticate(c *fiber.Ctx) error {
	s, err := h.store.Get(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if err = s.Destroy(); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}

// Get details of current user.
func (h *Handler) Get(c *fiber.Ctx) error {
	s, err := h.store.Get(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	userEmail := s.Get("email")
	if userEmail == nil {
		return fiber.ErrUnauthorized
	}

	user, err := h.backend.Get(userEmail.(string))
	if err != nil {
		log.Printf("Error getting current user : %v", err)

		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// Delete current user.
func (h *Handler) Delete(c *fiber.Ctx) error {
	s, err := h.store.Get(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	userEmail := s.Get("email")
	if userEmail == nil {
		return fiber.ErrUnauthorized
	}

	if err := h.backend.Delete(userEmail.(string)); err != nil {
		log.Printf("Error deleting user : %v", err)

		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}

// Verify user.
func (h *Handler) Verify(c *fiber.Ctx) error {
	s, err := h.store.Get(c)
	if err != nil || s.Get("email") == nil {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}

// Verify admin.
func (h *Handler) VerifyAdmin(c *fiber.Ctx) error {
	s, err := h.store.Get(c)
	if err != nil || s.Get("email") == nil || !s.Get("admin").(bool) {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}
