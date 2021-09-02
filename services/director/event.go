package main

import (
	"time"
)

// Event to gether user information on each redirect.
type Event struct {
	Time   time.Time
	Link   string
	Tag    string
	Cookie string
	UA     string
	IP     string
}

func NewEvent(link, tag, cookie, ua, ip string) Event {
	return Event{
		Time:   time.Now(),
		Link:   link,
		Tag:    tag,
		Cookie: cookie,
		UA:     ua,
		IP:     ip,
	}
}
