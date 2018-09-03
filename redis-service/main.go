package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.metrosystems.net/reliability-engineering/reliability-sandbox/GopherLab/redis-service/models"
	"git.metrosystems.net/reliability-engineering/reliability-sandbox/GopherLab/redis-service/utils"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	redisURI = "redis://:@localhost:6379/1"
	err      error
	client   *redis.Client
)

func init() {
	client = utils.NewRedisClient(redisURI)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", PingHandler).Methods("GET")
	router.HandleFunc("/set-key", SetKeyHandler).Methods("POST")
	router.HandleFunc("/get-key/{key}", GetKeyHandler).Methods("GET")

	handler := cors.AllowAll().Handler(router)

	if err := http.ListenAndServe("localhost:8080", handler); err != nil {
		log.Fatalf("ListenAndServe: %v", err.Error())
	}
}

func SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var params models.SetKeyParams
	utils.SafeParams(&params, r)
	err = client.Set(params.Key, params.Value, 0).Err()
	response := models.OutResponse{Message: "Success", Status: 200}
	if err != nil {
		response.Message = fmt.Sprintf("Error: %v", err)
		response.Status = 422
	}
	json.NewEncoder(w).Encode(response)
}

func GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	keyString := vars["key"]
	val, err := client.Get(keyString).Result()

	var result models.OutResponse
	if err == redis.Nil {
		result.Message = "Key not found"
		result.Status = 404
	} else if err != nil {
		result.Message = "Something went wrong"
		result.Status = 500
	} else {
		json.NewEncoder(w).Encode(models.SetKeyParams{Key: keyString, Value: val})
		return
	}

	json.NewEncoder(w).Encode(result)

}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	pong, err := client.Ping().Result()
	utils.CheckErr(err)
	json.NewEncoder(w).Encode(models.OutResponse{Message: pong, Status: 200})
}
