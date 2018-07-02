package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KeyValueList struct {
	Key  string   `json:"key"`
	List []string `json:"list"`
}

func createRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	redisClient := createRedisClient()

	// GET /ping
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/v1/value", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		// GET /v1/value
		case "GET":
			keyQuery := req.URL.Query()["key"]
			if len(keyQuery) != 1 {
				w.WriteHeader(400)
				return
			}

			value, err := redisClient.Get("Default::" + keyQuery[0]).Result()
			if err != nil && err != redis.Nil {
				w.WriteHeader(503)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(KeyValue{Key: keyQuery[0], Value: value})

		// POST /v1/value
		case "POST":
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				w.WriteHeader(500)
				return
			}

			var keyValue KeyValue
			err = json.Unmarshal(body, &keyValue)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			err = redisClient.Set("Default::"+keyValue.Key, keyValue.Value, 0).Err()
			if err != nil {
				w.WriteHeader(500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(KeyValue{Key: keyValue.Key, Value: keyValue.Value})
		}
	})

	http.HandleFunc("/v1/list", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		// GET /v1/list
		case "GET":
			keyQuery := req.URL.Query()["key"]
			if len(keyQuery) != 1 {
				w.WriteHeader(400)
				return
			}

			value, err := redisClient.LRange("Default::"+keyQuery[0], 0, -1).Result()
			if err != nil && err != redis.Nil {
				w.WriteHeader(503)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(KeyValueList{Key: keyQuery[0], List: value})

		// POST /v1/list
		case "POST":
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				w.WriteHeader(500)
				return
			}

			var keyValue KeyValue
			err = json.Unmarshal(body, &keyValue)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			err = redisClient.RPush("Default::"+keyValue.Key, keyValue.Value).Err()
			if err != nil {
				w.WriteHeader(500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(KeyValue{Key: keyValue.Key, Value: keyValue.Value})

		// DELETE /v1/list
		case "DELETE":
			keyQuery := req.URL.Query()["key"]
			if len(keyQuery) != 1 {
				w.WriteHeader(400)
				return
			}
			key := keyQuery[0]

			deleteKeyQuery := req.URL.Query()["delete_key"]
			valueQuery := req.URL.Query()["value"]

			if len(deleteKeyQuery) == 1 && deleteKeyQuery[0] == "true" {
				redisClient.Del("Default::" + key)
			} else if len(valueQuery) == 1 {
				redisClient.LRem("Default::"+key, 0, valueQuery[0])
			} else {
				w.WriteHeader(400)
				return
			}
		}
	})

	http.ListenAndServe(":8080", nil)
}
