package links

import (
	"context"
	"errors"
	"log"

	"github.com/go-redis/redis/v8"
)

// Redis implementation of link db strore

type RdbStore struct {
	db *redis.Client
}

func NewRdbStore(db *redis.Client) *RdbStore {
	return &RdbStore{
		db,
	}
}

func (r *RdbStore) Insert(link *Link) error {
	_, err := r.db.HSet(context.Background(), "links", link.Id, link.Target).Result()
	if err != nil {
		log.Printf("Error creating link : %v", err)
	}
	return err
}

func (r *RdbStore) Update(link *Link) error {
	_, err := r.db.HSet(context.Background(), "links", link.Id, link.Target).Result()
	if err != nil {
		log.Printf("Error updating link : %v", err)
	}
	return err
}

func (r *RdbStore) Get(id string) (*Link, error) {
	target, err := r.db.HGet(context.Background(), "links", id).Result()
	if err != nil {
		log.Printf("Link not found : %v", err)
		return nil, errors.New("Link not found")
	}
	return &Link{
		Id:     id,
		Target: target,
		Tag:    "None",
	}, nil
}

func (r *RdbStore) Delete(id string) error {
	_, err := r.db.HDel(context.Background(), "links", id).Result()
	if err != nil {
		log.Printf("Error deleting link %v", err)
	}
	return err
}

func (r *RdbStore) Ids() ([]string, error) {
	ids, err := r.db.HKeys(context.Background(), "links").Result()
	if err != nil {
		log.Printf("Error getting link ids %v", err)
		return nil, err
	}
	return ids, nil
}
