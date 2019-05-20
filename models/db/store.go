package db

import (
	"github.com/go-redis/redis"
)

type Store struct {
	connection *redis.Client
}

func NewStore() (*Store, error) {
	client, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	return &Store{
		connection: client,
	}, nil
}
