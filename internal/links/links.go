package links

// Link model and constructor

type Link struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Tag    string `json:"tag"`
}

func New(id, target, tag string) *Link {
	return &Link{
		Id:     id,
		Target: target,
		Tag:    tag,
	}
}
