package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func test(client *redis.Client, ctx context.Context) {
	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		log.Println("[!] Failed to set foo: %w", err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		log.Println("[!] Failed to get foo: %w", err)
	}
	fmt.Println("foo", val)
}

func getAllData(client *redis.Client, ctx context.Context) {
	// Get all keys and values
	keys, err := client.Keys(ctx, "*").Result()
	if err != nil {
		log.Println("[!] Failed to get keys: %w", err)
	}

	for _, key := range keys {
		val, err := client.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(key, val)
	}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	getAllData(client, ctx)
}
