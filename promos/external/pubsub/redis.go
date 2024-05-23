package pubsub

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(port string) *redis.Client {
    redisHost := fmt.Sprintf("localhost:%s", port)
    client := redis.NewClient(&redis.Options{
        Addr: redisHost,
        DB: 0,
    })
    return client
} 
