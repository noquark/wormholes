package ingestor

import (
	"context"
	"log"
	"time"
	"wormholes/services/creator/links"
	"wormholes/services/creator/sql"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	TickerInterval = time.Second * 10
)

// A simple link ingestor.
type Ingestor struct {
	db        *pgxpool.Pool
	batchSize int
	batch     *pgx.Batch
	quit      chan struct{}
	source    chan *links.Link
	ticker    *time.Ticker
}

func New(db *pgxpool.Pool, batchSize int) *Ingestor {
	return &Ingestor{
		db:        db,
		batchSize: batchSize,
		batch:     &pgx.Batch{},
		quit:      make(chan struct{}),
		source:    make(chan *links.Link),
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

func (i *Ingestor) Push(link *links.Link) {
	i.source <- link
}

func (i *Ingestor) add(link *links.Link) {
	i.batch.Queue(
		sql.Insert,
		link.ID, link.Tag, link.Target)

	if i.batch.Len() > i.batchSize {
		i.ingest()
	}
}

func (i *Ingestor) ingest() {
	batchOp := i.db.SendBatch(context.Background(), i.batch)

	_, err := batchOp.Exec()
	if err != nil {
		log.Printf("error inserting item : %v", err)
	}

	err = batchOp.Close()
	if err != nil {
		log.Printf("error closing batch : %v", err)
	}

	i.batch = &pgx.Batch{}
}
