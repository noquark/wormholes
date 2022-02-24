package director

import (
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const TickerInterval = time.Second * 10

// Pipe with multiple streams to ingest data concurrently
// For now Event is only data being ingested, but this
// implementation can handle more than one type of data

type Pipe struct {
	Streams []*Stream
	Task
	Queue
	db        *pgxpool.Pool
	batchSize int
	size      int
	ticker    *time.Ticker
}

func NewPipe(conf *Config, db *pgxpool.Pool) *Pipe {
	return &Pipe{
		Streams:   make([]*Stream, 0),
		Task:      make(Task),
		Queue:     make(Queue),
		db:        db,
		batchSize: conf.BatchSize,
		size:      conf.Streams,
		ticker:    time.NewTicker(TickerInterval),
	}
}

func (p *Pipe) Start() *Pipe {
	for i := 0; i < p.size; i++ {
		stream := NewStream(p.Queue, p.db, p.batchSize)
		go stream.Start()
		p.Streams = append(p.Streams, stream)
	}

	return p
}

func (p *Pipe) Wait() *Pipe {
	go func() {
		defer p.ticker.Stop()

		for {
			select {
			case item := <-p.Task:
				task := <-p.Queue
				task <- item
			case <-p.ticker.C:
				for _, s := range p.Streams {
					if s.Batch.Len() > 0 {
						s.Ingest()
					}
				}
			}
		}
	}()

	return p
}

func (p *Pipe) Push(item interface{}) {
	p.Task <- item
}

func (p *Pipe) Close() {
	for _, s := range p.Streams {
		if s.Batch.Len() > 0 {
			s.Ingest()
		}

		s.Stop()
	}

	log.Println("All streams are closed")
}
