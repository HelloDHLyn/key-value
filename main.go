package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	server.GET("/v1/value", func(c *gin.Context) {
		key := c.Query("key")
		value, err := redisClient.Get("Default::" + key).Result()
		if err != nil && err != redis.Nil {
			c.String(503, "")
			return
		}

		c.JSON(200, KeyValue{Key: key, Value: value})
	})

	server.POST("/v1/value", func(c *gin.Context) {
		var keyValue KeyValue
		c.ShouldBindJSON(&keyValue)

		err := redisClient.Set("Default::"+keyValue.Key, keyValue.Value, 0).Err()
		if err != nil {
			c.String(500, "")
			return
		}

		c.JSON(200, KeyValue{Key: keyValue.Key, Value: keyValue.Value})
	})

	server.GET("/v1/list", func(c *gin.Context) {
		key := c.Query("key")
		value, err := redisClient.LRange("Default::"+key, 0, -1).Result()
		if err != nil && err != redis.Nil {
			c.String(503, "")
			return
		}

		c.JSON(200, KeyValueList{Key: key, List: value})
	})

	server.POST("/v1/list", func(c *gin.Context) {
		var keyValue KeyValue
		c.ShouldBindJSON(&keyValue)

		err := redisClient.RPush("Default::"+keyValue.Key, keyValue.Value).Err()
		if err != nil {
			c.String(500, "")
			return
		}

		c.JSON(200, KeyValue{Key: keyValue.Key, Value: keyValue.Value})
	})

	server.DELETE("/v1/list", func(c *gin.Context) {
		key := c.Query("key")
		if len(key) == 0 {
			c.String(400, "")
			return
		}

		deleteKey := c.Query("delete_key")
		value := c.Query("value")

		if deleteKey == "true" {
			redisClient.Del("Default::" + key)
			c.String(200, "")
		} else if value != "" {
			redisClient.LRem("Default::"+key, 0, value)
			c.String(200, "")
		} else {
			c.String(400, "")
			return
		}
	})

	server.Run(":8080")
}
