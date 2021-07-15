package pipe

import "github.com/jackc/pgx/v4/pgxpool"

type Task chan Event
type Queue chan chan Event

type Stream struct {
	Task
	Queue
	Quit chan struct{}
	db   *pgxpool.Pool
}

func NewStream(queue Queue, db *pgxpool.Pool) *Stream {
	stream := &Stream{
		Task:  make(Task),
		Queue: queue,
		Quit:  make(chan struct{}),
	}
	return stream
}

func (s *Stream) Start() {
	for {
		s.Queue <- s.Task
		select {
		case event := <-s.Task:
			s.Ingest(&event)
		case <-s.Quit:
			close(s.Task)
			return
		}
	}
}

func (s *Stream) Stop() {
	close(s.Quit)
}

func (s *Stream) Ingest(event *Event) {
	// TODO: Insert event to timescale
}
