package pipe

import (
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mohitsinghs/wormholes/config"
	"github.com/oschwald/geoip2-golang"
)

// Pipe with multiple streams to ingest data concurrently

type Pipe struct {
	Streams []*Stream
	Task
	Queue
	db        *pgxpool.Pool
	ip        *geoip2.Reader
	batchSize int
	size      int
	ticker    *time.Ticker
}

func New(conf *config.Config) *Pipe {
	return &Pipe{
		Streams:   make([]*Stream, 0),
		Task:      make(Task),
		Queue:     make(Queue),
		db:        conf.Postgres.Connect(),
		ip:        conf.Pipe.OpenDB(),
		batchSize: conf.Pipe.BatchSize,
		size:      conf.Pipe.Streams,
		ticker:    time.NewTicker(10 * time.Second),
	}
}

func (p *Pipe) Start() *Pipe {
	for i := 0; i < p.size; i++ {
		stream := NewStream(p.Queue, p.db, p.ip, p.batchSize)
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
