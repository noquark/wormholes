package redirector

import (
	"context"
	_ "embed"
	"reflect"
	"time"
	"wormholes/internal/cache"
	"wormholes/internal/links"
	"wormholes/services/redirector/pipe"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	CookieExpiryTime = time.Hour * 24 * 180
	CacheControl     = "private, max-age=90"
	CookieName       = "_wh"
)

//go:embed sql/get_link.sql
var linkGet string

// Fiber route handlers for link

type Handler struct {
	pipe  *pipe.Pipe
	db    *pgxpool.Pool
	cache *cache.Cache
}

func NewHandler(pipe *pipe.Pipe, db *pgxpool.Pool, cache *cache.Cache) *Handler {
	return &Handler{
		pipe,
		db,
		cache,
	}
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
		err := h.db.QueryRow(context.Background(),
			linkGet,
			shortID,
		).Scan(&link.ID, &link.Target, &link.Tag)
		if err != nil {
			log.Error().Err(err).Msg("redirect: error getting link")

			return fiber.ErrInternalServerError
		}

		err = h.cache.SetLink(link, shortID)
		if err != nil {
			log.Warn().Err(err).Msg("redirect: failed to cache")
		}
	}

	var cookie string

	if c.Cookies(CookieName) == "" {
		cookie := NewCookie()

		c.Cookie(&fiber.Cookie{
			Name:    CookieName,
			Value:   cookie,
			Expires: time.Now().Add(CookieExpiryTime),
		})
	} else {
		cookie = c.Cookies(CookieName)
	}

	c.Set(fiber.HeaderCacheControl, CacheControl)
	h.pipe.Push(pipe.NewEvent(link.ID, link.Tag, cookie, c.Get(fiber.HeaderUserAgent), c.IP()))

	return c.Redirect(link.Target, fiber.StatusMovedPermanently)
}
