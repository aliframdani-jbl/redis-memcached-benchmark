package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
)

func redisSet(client *redis.Client, key, value string) error {
	return client.Set(ctx, key, value, 0).Err()
}

func redisGet(client *redis.Client, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func BenchmarkRedisStringSetGet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	value := "true"

	const numGoroutines = 100
	const numIterations = 100000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	errs := make(chan error, numGoroutines)

	b.ResetTimer()

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations/numGoroutines; j++ {
				key := randMarketplace() + ":" + randMarketplaceAndStoreId()
				if err := redisSet(client, key, value); err != nil {
					errs <- err
					return
				}
				retrievedValue, err := redisGet(client, key)
				if err != nil {
					errs <- err
					return
				}
				if retrievedValue != value {
					errs <- fmt.Errorf("expected %s, got %s", value, retrievedValue)
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			b.Fatal(err)
		}
	}
}
