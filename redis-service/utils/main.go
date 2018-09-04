package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"git.metrosystems.net/reliability-engineering/reliability-sandbox/GopherLab/redis-service/models"
	"github.com/go-redis/redis"
)

// NewRedisClient will parse a redis connection string and will return a redis client
func NewRedisClient(redisURI string) *redis.Client {
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}

// CheckErr will print to log an existing error
func CheckErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

// SafeParams is a helper that filters and extracts only desired params from request body
func SafeParams(params interface{}, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	CheckErr(err)
}

func SerializeErrMessage(response models.OutResponse) string {
	raw, err := json.Marshal(response)
	CheckErr(err)
	return string(raw)
}
