package pipe

import (
	"context"
	_ "embed"
	"log"
	"wormholes/services/redirector/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mssola/user_agent"
)

// Stream to ingest data in batches

type Stream struct {
	Task
	Queue
	Quit      chan struct{}
	db        *pgxpool.Pool
	batchSize int
	Batch     *pgx.Batch
}

func NewStream(queue Queue, db *pgxpool.Pool, batchSize int) *Stream {
	stream := &Stream{
		Task:      make(Task),
		Queue:     queue,
		Quit:      make(chan struct{}),
		db:        db,
		batchSize: batchSize,
		Batch:     &pgx.Batch{},
	}

	return stream
}

func (s *Stream) Start() {
	for {
		s.Queue <- s.Task
		select {
		case item := <-s.Task:
			s.Add(item)
		case <-s.Quit:
			close(s.Task)

			return
		}
	}
}

func (s *Stream) Stop() {
	close(s.Quit)
}

func (s *Stream) Add(item interface{}) {
	switch item := item.(type) {
	case Event:
		ua := user_agent.New(item.UA)
		browser, browserVersion := ua.Browser()

		s.Batch.Queue(
			sql.InsertEvent,
			item.Time, item.Link, item.Tag, item.Cookie, item.IP,
			ua.Mobile(), ua.Bot(), browser, browserVersion, ua.OSInfo().Name, ua.OSInfo().Version, ua.Platform(),
		)
	default:
		log.Println("ignoring item")
	}

	if s.Batch.Len() > s.batchSize {
		s.Ingest()
	}
}

func (s *Stream) Ingest() {
	br := s.db.SendBatch(context.Background(), s.Batch)

	_, err := br.Exec()
	if err != nil {
		log.Printf("error inserting item : %v", err)
	}

	err = br.Close()
	if err != nil {
		log.Printf("error closing batch : %v", err)
	}

	s.Batch = &pgx.Batch{}
}
