package main

import (
	"context"
	_ "embed"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
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
	pipe  *Pipe
	db    *pgxpool.Pool
	cache *redis.Client
}

func NewHandler(pipe *Pipe, db *pgxpool.Pool, cache *redis.Client) *Handler {
	return &Handler{
		pipe,
		db,
		cache,
	}
}

func (h *Handler) Redirect(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return fiber.ErrBadRequest
	}

	var link Link

	err := h.cache.HGetAll(context.Background(), id).Scan(&link)
	if err != nil {
		log.Err(err).Msg("redirect: cache miss")

		// If key does not exists, query db
		rows := h.db.QueryRow(context.Background(),
			linkGet,
			id,
		)

		err := rows.Scan(&link.ID, &link.Target, &link.Tag)
		if err != nil {
			log.Error().Err(err).Msg("redirect: error getting link")

			return fiber.ErrInternalServerError
		}

		_ = h.cache.HSet(context.Background(), id, link).Err()
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
	h.pipe.Push(NewEvent(link.ID, link.Tag, cookie, c.Get(fiber.HeaderUserAgent), c.IP()))

	return c.Redirect(link.Target, fiber.StatusMovedPermanently)
}
