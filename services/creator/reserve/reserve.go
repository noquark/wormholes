package reserve

type Reserve interface {
	GetID() (string, error)
}
