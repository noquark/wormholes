package links

// Link model and constructor

type Link struct {
	ID     string `json:"id" redis:"id"`
	Target string `json:"target" redis:"target"`
	Tag    string `json:"tag" redis:"tag"`
}

func New(id, target, tag string) *Link {
	return &Link{
		ID:     id,
		Target: target,
		Tag:    tag,
	}
}
