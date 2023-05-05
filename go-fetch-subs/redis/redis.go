package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"

	"github.com/EC3-Gang/cbr-api/scraper"
)

type redisClient struct {
	client *redis.Client
	ctx    context.Context
}

func test(r redisClient) {
	err := r.client.Set(r.ctx, "foo", "bar", 0).Err()
	if err != nil {
		log.Println("[!] Failed to set foo: %w", err)
	}

	val, err := r.client.Get(r.ctx, "foo").Result()
	if err != nil {
		log.Println("[!] Failed to get foo: %w", err)
	}
	fmt.Println("foo", val)
}

func getAllData(r redisClient) {
	keys, err := r.client.Keys(r.ctx, "*").Result()
	if err != nil {
		log.Println("[!] Failed to get keys: %w", err)
	}

	for _, key := range keys {
		val, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(key, val)
	}
}

func storeAttempts(r redisClient, name string, attempts []scraper.Attempt) {
	err := r.client.Set(r.ctx, name, attempts, 0).Err()
	if err != nil {
		log.Println("[!] Failed to set attempts: %w", err)
	}
}

func getAttempts(r redisClient, name string) []scraper.Attempt {
	val, err := r.client.Get(r.ctx, name).Result()
	if err != nil {
		log.Println("[!] Failed to get attempts: %w", err)
	}
	// Convert string to []scraper.Attempt
	var attempts []scraper.Attempt
	err = json.Unmarshal([]byte(val), &attempts)
	if err != nil {
		log.Println("[!] Failed to unmarshal attempts: %w", err)
	}
	return attempts
}

func newClient(host string, port int) redisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	return redisClient{client: client, ctx: ctx}
}
