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

type Task chan Event
type Queue chan chan Event

const sqlInsert = `insert into wh_clicks (time,link,tag,cookie,ip,is_mobile,is_bot,browser,browser_version,os, os_version,platform,lat,long,city,region,country,continent) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18);`

type Stream struct {
	Task
	Queue
	Quit      chan struct{}
	db        *pgxpool.Pool
	ip        *geoip2.Reader
	batchSize int
	batch     *pgx.Batch
}

func NewStream(queue Queue, db *pgxpool.Pool, ip *geoip2.Reader, batchSize int) *Stream {
	stream := &Stream{
		Task:      make(Task),
		Queue:     queue,
		Quit:      make(chan struct{}),
		db:        db,
		ip:        ip,
		batchSize: batchSize,
		batch:     &pgx.Batch{},
	}
	return stream
}

func (s *Stream) Start() {
	for {
		s.Queue <- s.Task
		select {
		case event := <-s.Task:
			s.Add(&event)
		case <-s.Quit:
			close(s.Task)
			return
		}
	}
}

func (s *Stream) Stop() {
	close(s.Quit)
}

func (s *Stream) Add(event *Event) {
	ua := user_agent.New(event.UA)
	browser, browser_version := ua.Browser()
	address, _ := s.ip.City(net.ParseIP(event.IP))
	region := ""
	if len(address.Subdivisions) > 0 {
		region = address.Subdivisions[0].Names[constants.EN]
	}
	s.batch.Queue(
		sqlInsert,
		event.Time, event.Link, event.Tag, event.Cookie, event.IP,
		ua.Mobile(), ua.Bot(), browser, browser_version, ua.OSInfo().Name, ua.OSInfo().Version, ua.Platform(),
		address.Location.Latitude, address.Location.Longitude, address.City.Names[constants.EN], region, address.Country.Names[constants.EN], address.Continent.Names[constants.EN],
	)
	if s.batch.Len() > s.batchSize {
		s.Ingest()
	}
}

func (s *Stream) Ingest() {
	br := s.db.SendBatch(context.Background(), s.batch)
	_, err := br.Exec()
	if err != nil {
		log.Printf("error inserting event : %v", err)
	}
	err = br.Close()
	if err != nil {
		log.Printf("error closing batch : %v", err)
	}
	s.batch = &pgx.Batch{}
}
