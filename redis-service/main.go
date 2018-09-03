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
	log.Println("Server listening on localhost:8080")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", PingHandler).Methods("GET")
	router.HandleFunc("/set-key", SetKeyHandler).Methods("POST")
	router.HandleFunc("/get-key/{key}", GetKeyHandler).Methods("GET")

	handler := cors.AllowAll().Handler(router)

	if err := http.ListenAndServe("localhost:8080", handler); err != nil {
		log.Fatalf("ListenAndServe: %v", err.Error())
	}
}

// SetKeyHandler is used to set a value for a specified key
func SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var params models.SetKeyParams
	utils.SafeParams(&params, r)
	log.Printf("<Set key> params: %v\n", params)
	err = client.Set(params.Key, params.Value, 0).Err()
	response := models.OutResponse{Message: "Success", Status: http.StatusOK}
	if err != nil {
		response.Message = fmt.Sprintf("Error: %v", err)
		response.Status = http.StatusUnprocessableEntity
		log.Printf("<Set key> error: %v\n", err)
	}

	log.Printf("<Set key> result: %v\n", response)
	json.NewEncoder(w).Encode(response)
}

// GetKeyHandler is used to fetch a value for a specified key
func GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	keyString := vars["key"]
	log.Printf("<Get key> key: %v\n", keyString)
	val, err := client.Get(keyString).Result()

	var result models.OutResponse
	if err == redis.Nil {
		result.Message = "Key not found"
		result.Status = http.StatusNotFound
	} else if err != nil {
		result.Message = "Something went wrong"
		result.Status = http.StatusInternalServerError
		log.Printf("<Get key> error: %v\n", err)
	} else {
		json.NewEncoder(w).Encode(models.SetKeyParams{Key: keyString, Value: val})
		return
	}

	json.NewEncoder(w).Encode(result)

}

// PingHandler is used to make sure connection to redis server is ok
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	pong, err := client.Ping().Result()
	utils.CheckErr(err)
	json.NewEncoder(w).Encode(models.OutResponse{Message: pong, Status: http.StatusOK})
}
