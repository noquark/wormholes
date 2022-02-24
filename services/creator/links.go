package creator

// Link model and constructor

type Link struct {
	ID     string `json:"id"`
	Target string `json:"target"`
	Tag    string `json:"tag"`
}

func NewLink(id, target, tag string) *Link {
	return &Link{
		ID:     id,
		Target: target,
		Tag:    tag,
	}
}
