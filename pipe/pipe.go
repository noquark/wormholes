package pipe

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mohitsinghs/wormholes/config"
)

type Pipe struct {
	Streams []*Stream
	Task
	Queue
	db *pgxpool.Pool
}

func New(conf *config.Config) *Pipe {
	pipe := &Pipe{
		Streams: make([]*Stream, conf.Streams),
		Task:    make(Task),
		Queue:   make(Queue),
		db:      conf.Postgres.Connect(),
	}
	return pipe
}

func (p *Pipe) Start() *Pipe {
	size := len(p.Streams)
	for i := 0; i < size; i++ {
		stream := NewStream(p.Queue, p.db)
		go stream.Start()
		p.Streams = append(p.Streams, stream)
	}
	return p
}

func (p *Pipe) Wait() *Pipe {
	go func() {
		for {
			event := <-p.Task
			task := <-p.Queue
			task <- event
		}
	}()
	return p
}

func (p *Pipe) Push(event Event) {
	p.Task <- event
}
