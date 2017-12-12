package socketstore

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var client *redis.Client

func setupRedis() {
	redisUrl := os.Getenv("LOCATION_REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}
	redisPassword := os.Getenv("LOCATION_REDIS_PASSWORD")
	client = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       2,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func GetClient() *redis.Client {
	return client
}
