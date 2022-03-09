package creator

type Store interface {
	Get(id string) (Link, error)
	Update(link *Link) error
	Delete(id string) error
}
