package creator

import (
	"context"
	_ "embed"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed sql/insert_link.sql
var linkInsert string

const (
	TickerInterval = time.Second * 10
)

// A simple link ingestor.
type Ingestor struct {
	db        *pgxpool.Pool
	batchSize int
	batch     *pgx.Batch
	quit      chan struct{}
	source    chan *Link
	ticker    *time.Ticker
}

func NewIngestor(db *pgxpool.Pool, batchSize int) *Ingestor {
	return &Ingestor{
		db:        db,
		batchSize: batchSize,
		batch:     &pgx.Batch{},
		quit:      make(chan struct{}),
		source:    make(chan *Link),
		ticker:    time.NewTicker(TickerInterval),
	}
}

func (i *Ingestor) Start() *Ingestor {
	go func() {
		defer i.ticker.Stop()

		for {
			select {
			case link := <-i.source:
				i.add(link)
			case <-i.quit:
				close(i.source)

				return
			case <-i.ticker.C:
				if i.batch.Len() > 0 {
					i.ingest()
				}
			}
		}
	}()

	return i
}

func (i *Ingestor) Push(link *Link) {
	i.source <- link
}

func (i *Ingestor) add(link *Link) {
	i.batch.Queue(
		linkInsert,
		link.ID, link.Tag, link.Target)

	if i.batch.Len() > i.batchSize {
		i.ingest()
	}
}

func (i *Ingestor) ingest() {
	br := i.db.SendBatch(context.Background(), i.batch)

	_, err := br.Exec()
	if err != nil {
		log.Printf("error inserting item : %v", err)
	}

	err = br.Close()
	if err != nil {
		log.Printf("error closing batch : %v", err)
	}

	i.batch = &pgx.Batch{}
}
