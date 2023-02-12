package redisUtil

import (
	"github.com/go-redis/redis"
)

var Rdb *redis.Client

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root", // no password set
		DB:       0,      // use default DB
	})
}

//	pong, err := rdb.Ping().Result()
//	fmt.Println(pong, err)
