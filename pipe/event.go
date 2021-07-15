package pipe

import (
	"time"
)

type Event struct {
	Time   time.Time
	Link   string
	Tag    string
	Cookie string
	UA     string
	IP     string
}
