package stats

type Overview struct {
	Links  uint64 `json:"links"`
	Tags   uint64 `json:"tags"`
	Clicks uint64 `json:"clicks"`
	Users  uint64 `json:"users"`
}
