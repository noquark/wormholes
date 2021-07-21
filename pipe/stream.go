package pipe

import (
	"context"
	_ "embed"
	"log"
	"net"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mohitsinghs/wormholes/constants"
	"github.com/mssola/user_agent"
	"github.com/oschwald/geoip2-golang"
)

//go:embed sql/insert_event.sql
var eventInsert string

// Stream to ingest data in batches

type Stream struct {
	Task
	Queue
	Quit      chan struct{}
	db        *pgxpool.Pool
	ip        *geoip2.Reader
	batchSize int
	Batch     *pgx.Batch
}

func NewStream(queue Queue, db *pgxpool.Pool, ip *geoip2.Reader, batchSize int) *Stream {
	stream := &Stream{
		Task:      make(Task),
		Queue:     queue,
		Quit:      make(chan struct{}),
		db:        db,
		ip:        ip,
		batchSize: batchSize,
		Batch:     &pgx.Batch{},
	}
	return stream
}

func (s *Stream) Start() {
	for {
		s.Queue <- s.Task
		select {
		case event := <-s.Task:
			s.Add(event)
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
		browser, browser_version := ua.Browser()
		address, _ := s.ip.City(net.ParseIP(item.IP))
		region := ""
		if len(address.Subdivisions) > 0 {
			region = address.Subdivisions[0].Names[constants.EN]
		}
		s.Batch.Queue(
			eventInsert,
			item.Time, item.Link, item.Tag, item.Cookie, item.IP,
			ua.Mobile(), ua.Bot(), browser, browser_version, ua.OSInfo().Name, ua.OSInfo().Version, ua.Platform(),
			address.Location.Latitude, address.Location.Longitude, address.City.Names[constants.EN], region, address.Country.Names[constants.EN], address.Continent.Names[constants.EN],
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
