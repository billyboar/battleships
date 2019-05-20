package db

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

// ConnectDB connects to local redis instance
func ConnectDB() (*redis.Client, error) {
	fmt.Println("DEBUG REDIS!!!!!", os.Getenv("REDIS_HOST"))
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
