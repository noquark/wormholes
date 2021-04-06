package links

// Interface for link db store
// replacable bridge between db and actual handler

type Store interface {
	Insert(*Link) error
	Update(*Link) error
	Get(id string) (*Link, error)
	Delete(id string) error
	Ids() ([]string, error)
}
