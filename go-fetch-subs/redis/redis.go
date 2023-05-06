package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EC3-Gang/cbr-api/types"
	"github.com/redis/go-redis/v9"
	"log"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
}

func test(r Client) {
	err := r.client.Set(r.ctx, "foo", "bar", 0).Err()
	if err != nil {
		log.Printf("[!] Failed to set foo: %v", err)
	}

	val, err := r.client.Get(r.ctx, "foo").Result()
	if err != nil {
		log.Printf("[!] Failed to get foo: %v", err)
	}
	fmt.Println("foo", val)
}

func getAllData(r Client) {
	keys, err := r.client.Keys(r.ctx, "*").Result()
	if err != nil {
		log.Printf("[!] Failed to get keys: %v", err)
	}

	for _, key := range keys {
		val, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(key, val)
	}
}

func storeAttempts(r Client, name string, attempts *types.AttemptList) {
	marshalledAttempts, err := attempts.MarshalBinary()
	if err != nil {
		log.Printf("[!] Failed to marshal attempts: %v", err)
	}

	err = r.client.Set(r.ctx, name, marshalledAttempts, 0).Err()
	if err != nil {
		log.Printf("[!] Failed to set attempts: %v", err)
	}
}

func getAttempts(r Client, name string) *[]types.Attempt {
	val, err := r.client.Get(r.ctx, name).Result()
	if err != nil {
		log.Printf("[!] Failed to get attempts: %v", err)
	}
	// Convert string to []scraper.Attempt
	var attempts []types.Attempt
	err = json.Unmarshal([]byte(val), &attempts)
	if err != nil {
		log.Printf("[!] Failed to unmarshal attempts: %v", err)
	}
	return &attempts
}

func addProblem(r Client, problemID string) {
	err := r.client.SAdd(r.ctx, "problems", problemID).Err()
	if err != nil {
		log.Printf("[!] Failed to add problem: %v", err)
	}
}

func checkProblemCached(r Client, problemID string) bool {
	fmt.Println(r.client.SIsMember(r.ctx, "problems", problemID).Val())
	return r.client.SIsMember(r.ctx, "problems", problemID).Val()
}

func NewClient(host string, port int) Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	return Client{client: client, ctx: ctx}
}
