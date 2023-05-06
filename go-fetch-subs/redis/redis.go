package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EC3-Gang/cbr-api/types"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func test(r RedisClient) {
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

func getAllData(r RedisClient) {
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

func storeAttempts(r RedisClient, name string, attempts *[]types.Attempt) {
	err := r.client.Set(r.ctx, name, *attempts, 0).Err()
	if err != nil {
		log.Println("[!] Failed to set attempts: %w", err)
	}
}

func getAttempts(r RedisClient, name string) *[]types.Attempt {
	val, err := r.client.Get(r.ctx, name).Result()
	if err != nil {
		log.Println("[!] Failed to get attempts: %w", err)
	}
	// Convert string to []scraper.Attempt
	var attempts []types.Attempt
	err = json.Unmarshal([]byte(val), &attempts)
	if err != nil {
		log.Println("[!] Failed to unmarshal attempts: %w", err)
	}
	return &attempts
}

func addProblem(r RedisClient, problemID string) {
	err := r.client.SAdd(r.ctx, "problems", problemID).Err()
	if err != nil {
		log.Println("[!] Failed to add problem: %w", err)
	}
}

func checkProblemCached(r RedisClient, problemID string) bool {
	return r.client.SIsMember(r.ctx, "problems", problemID).Val()
}

func NewClient(host string, port int) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	return RedisClient{client: client, ctx: ctx}
}
