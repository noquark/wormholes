package ipc

import (
	"context"
	"errors"
	"sync"
	"time"
	"wormholes/protos"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// backOff time is 500ms by default.
	backOffTime = time.Millisecond * 250
	ErrNoIds    = errors.New("reserve: there are no IDs ready yet")
)

type Store struct {
	mutex  sync.RWMutex
	status *Status
	bucket *protos.Bucket
	conn   *grpc.ClientConn
	client protos.BucketServiceClient
}

func NewStore(port string) *Store {

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("grpc-reserve: grpc failed to connect")
	}

	client := protos.NewBucketServiceClient(conn)

	return &Store{
		mutex:  sync.RWMutex{},
		status: NewStatus(),
		bucket: &protos.Bucket{},
		conn:   conn,
		client: client,
	}
}

func (s *Store) isEmpty() bool {
	return len(s.bucket.Ids) == 0
}

func (s *Store) fetch() {
	if s.status.IsBusy() {
		return
	}

	s.status.SetBusy()
	defer s.status.SetIdle()

	bucket, err := s.client.GetBucket(context.Background(), &protos.Empty{})
	if err != nil {
		log.Error().Err(err).Msg("grpc-reserve: grpc failed to fetch bucket")
	}

	if len(bucket.Ids) > 0 {
		s.mutex.Lock()
		s.bucket.Ids = bucket.Ids
		s.mutex.Unlock()
	}
}

func (s *Store) pop() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := s.bucket.Ids[0]
	s.bucket.Ids = s.bucket.Ids[1:]

	return id
}

func (s *Store) GetID() (string, error) {
	if s.isEmpty() {
		s.fetch()
		// this delays some request instead of failing them
		time.Sleep(backOffTime)
	} else {
		return s.pop(), nil
	}

	if !s.isEmpty() {
		return s.pop(), nil
	}

	return "", ErrNoIds
}
