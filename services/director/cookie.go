package director

import "github.com/wormholesdev/nanoid"

const (
	MaxTry     = 10
	CookieSize = 21
)

// Generate a random cookie with retry on failure.
func NewCookie() string {
	cookie, err := nanoid.New(CookieSize)
	if err != nil {
		for i := 0; i < MaxTry; i++ {
			cookie, err = nanoid.New(CookieSize)
			if err != nil {
				continue
			}

			break
		}
	}

	return cookie
}
