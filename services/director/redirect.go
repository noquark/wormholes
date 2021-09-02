package main

import (
	"context"
	_ "embed"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
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
	cache *redis.Pool
}

func NewHandler(pipe *Pipe, db *pgxpool.Pool, cache *redis.Pool) *Handler {
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

	// retrieve data from cache first
	conn := h.cache.Get()
	val, err := redis.Values(conn.Do("HGETALL", id))
	errScan := redis.ScanStruct(val, &link)

	if err != nil || errScan != nil {
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

		_, _ = conn.Do("HSET", redis.Args{}.Add(id).AddFlat(link)...)
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
