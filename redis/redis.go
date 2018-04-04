package redis

import (
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var client *redis.Client

func SetupRedis() *redis.Client {
	redisUrl := os.Getenv("LOCATION_REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}
	redisPassword := os.Getenv("LOCATION_REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       2,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Errorf("Unable to connect to redis at %s - %s", redisUrl, err)
		panic(err)
	}
	log.Infof("Connection to redis at %s - established", redisUrl)
	log.Infof("Pong result is %s ", pong)
	return client
}

func SetClient(cl *redis.Client) {
	client = cl
}

func GetClient() *redis.Client {
	return client
}
