package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	router.HandleFunc("/del-keys", DelKeyHandler).Methods("DELETE")

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

	ttl := time.Duration(params.TTL) * time.Second
	if ttl == 0 {
		ttl = time.Duration(864000) * time.Second
	}
	err = client.Set(params.Key, params.Value, ttl).Err()
	response := models.OutResponse{Message: "Success", Status: http.StatusOK}
	if err != nil {
		response.Message = fmt.Sprintf("Error: %v", err)
		response.Status = http.StatusUnprocessableEntity
		log.Printf("<Set key> error: %v\n", err)
		http.Error(w, utils.SerializeErrMessage(response), http.StatusUnprocessableEntity)
		return
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
		http.Error(w, utils.SerializeErrMessage(result), http.StatusInternalServerError)
		return
	} else {
		ttl, _ := client.TTL(keyString).Result()
		json.NewEncoder(w).Encode(models.SetKeyParams{Key: keyString, Value: val, TTL: int(ttl / time.Second)})
		return
	}

	json.NewEncoder(w).Encode(result)

}

// PingHandler is used to make sure connection to redis server is ok
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	pong, err := client.Ping().Result()
	if err != nil {
		http.Error(w, "Service Unavailable", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(models.OutResponse{Message: pong, Status: http.StatusOK})
}

// DelKeyHandler is used to delete a list of keys
func DelKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := r.URL.Query()
	keys := vars["keys"]

	response := models.OutResponse{}
	if len(keys) == 0 {
		response.Message = "No keys provided"
		response.Status = http.StatusUnprocessableEntity
		log.Println("<Del keys> error: No keys provided")
		http.Error(w, utils.SerializeErrMessage(response), http.StatusUnprocessableEntity)
		return
	}
	err := client.Del(keys...).Err()

	if err != nil {
		response.Message = "Something went wrong"
		response.Status = http.StatusInternalServerError
		log.Printf("<Del keys> error: %v\n", err)
		http.Error(w, utils.SerializeErrMessage(response), http.StatusInternalServerError)
		return
	}
	response.Message = fmt.Sprintf("Successfully deleted keys: %v", keys)
	response.Status = http.StatusOK
	log.Printf("<Del keys> success: Successfully deleted keys: %v", keys)

	json.NewEncoder(w).Encode(response)
}
