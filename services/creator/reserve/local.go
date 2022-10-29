package reserve

import (
	"sync"
	"time"
	"wormholes/services/generator"

	"github.com/rs/zerolog/log"
)

type LocalReserve struct {
	mutex   sync.RWMutex
	status  Status
	factory *generator.Factory
	bucket  []string
}

func WithLocal(f *generator.Factory) *LocalReserve {
	return &LocalReserve{
		mutex:   sync.RWMutex{},
		status:  *NewStatus(),
		factory: f,
	}
}

func (r *LocalReserve) isEmpty() bool {
	return len(r.bucket) == 0
}

func (r *LocalReserve) fetch() {
	if r.status.IsBusy() {
		return
	}

	r.status.SetBusy()
	defer r.status.SetIdle()

	ids, err := r.factory.GetLocalBucket()
	if len(ids) == 0 || err != nil {
		log.Error().Err(err).Msg("local-reserve: failed to fetch bucket")
	}

	if len(ids) > 0 {
		r.mutex.Lock()
		r.bucket = ids
		r.mutex.Unlock()
	}
}

func (r *LocalReserve) pop() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	id := r.bucket[0]
	r.bucket = r.bucket[1:]

	return id
}

func (r *LocalReserve) GetID() (string, error) {
	if r.isEmpty() {
		r.fetch()
		time.Sleep(backOffTime)
	} else {
		return r.pop(), nil
	}

	if !r.isEmpty() {
		return r.pop(), nil
	}

	return "", ErrNoIds
}
