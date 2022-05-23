package reserve

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
	// backoff time is 500ms by default.
	backoffTime = time.Millisecond * 5e2
	ErrNoIds    = errors.New("grcp-reserve: there are no IDs ready yet")
)

type GrpcReserve struct {
	mutex  sync.RWMutex
	status Status
	bucket *protos.Bucket
	conn   *grpc.ClientConn
	client protos.BucketServiceClient
}

func WithGrpc(addr string) *GrpcReserve {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("grpc-reserve: grpc failed to connect")
	}

	client := protos.NewBucketServiceClient(conn)

	return &GrpcReserve{
		mutex:  sync.RWMutex{},
		status: *NewStatus(),
		bucket: &protos.Bucket{},
		conn:   conn,
		client: client,
	}
}

func (r *GrpcReserve) isEmpty() bool {
	return len(r.bucket.Ids) == 0
}

func (r *GrpcReserve) fetch() {
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

func (r *GrpcReserve) pop() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	id := r.bucket.Ids[0]
	r.bucket.Ids = r.bucket.Ids[1:]

	return id
}

func (r *GrpcReserve) GetID() (string, error) {
	if r.isEmpty() {
		r.fetch()
		// this delays some request instead of failing them
		time.Sleep(backoffTime)
	} else {
		return r.pop(), nil
	}

	if !r.isEmpty() {
		return r.pop(), nil
	}

	return "", ErrNoIds
}
