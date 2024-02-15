package reserve

import (
	"context"
	"errors"
	"protos"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// backOff time is 500ms by default.
	backOffTime = time.Millisecond * 5e2
	ErrNoIds    = errors.New("reserve: there are no IDs ready yet")
)

type Reserve struct {
	mutex  sync.RWMutex
	status Status
	bucket *protos.Bucket
	conn   *grpc.ClientConn
	client protos.BucketServiceClient
}

func WithGrpc(addr string) *Reserve {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("grpc-reserve: grpc failed to connect")
	}

	client := protos.NewBucketServiceClient(conn)

	return &Reserve{
		mutex:  sync.RWMutex{},
		status: *NewStatus(),
		bucket: &protos.Bucket{},
		conn:   conn,
		client: client,
	}
}

func (r *Reserve) isEmpty() bool {
	return len(r.bucket.Ids) == 0
}

func (r *Reserve) fetch() {
	if r.status.IsBusy() {
		return
	}

	r.status.SetBusy()
	defer r.status.SetIdle()

	bucket, err := r.client.GetBucket(context.Background(), &protos.Empty{})
	if err != nil {
		log.Error().Err(err).Msg("grpc-reserve: grpc failed to fetch bucket")
	}

	if len(bucket.Ids) > 0 {
		r.mutex.Lock()
		r.bucket.Ids = bucket.Ids
		r.mutex.Unlock()
	}
}

func (r *Reserve) pop() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	id := r.bucket.Ids[0]
	r.bucket.Ids = r.bucket.Ids[1:]

	return id
}

func (r *Reserve) GetID() (string, error) {
	if r.isEmpty() {
		r.fetch()
		// this delays some request instead of failing them
		time.Sleep(backOffTime)
	} else {
		return r.pop(), nil
	}

	if !r.isEmpty() {
		return r.pop(), nil
	}

	return "", ErrNoIds
}
