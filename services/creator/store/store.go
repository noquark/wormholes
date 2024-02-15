package store

import "lib/links"

type Store interface {
	Get(id string) (links.Link, error)
	Update(link *links.Link) error
	Delete(id string) error
}
