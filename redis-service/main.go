package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var (
	redisURI = "redis://:@localhost:6379/1"
	err      error
)

func newRedisClient() *redis.Client {
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	client := newRedisClient()
	pong, err := client.Ping().Result()
	checkErr(err)
	fmt.Println(pong)

	err = client.Set("key", "value", 0).Err()
	checkErr(err)

	val, err := client.Get("key").Result()
	checkErr(err)
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
