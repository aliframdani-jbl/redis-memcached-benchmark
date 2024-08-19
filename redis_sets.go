package main

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func redisSetMembers(client *redis.Client, setName string, members []string) error {
	_, err := client.SAdd(ctx, setName, members).Result()
	return err
}

func redisIsMember(client *redis.Client, setName, member string) (bool, error) {
	isMember, err := client.SIsMember(ctx, setName, member).Result()
	return isMember, err
}

func BenchmarkRedisSetMembers(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	setName := "myset"
	members := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		members[i] = fmt.Sprintf("member%d", i)
	}

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
				if err := redisSetMembers(client, setName, members); err != nil {
					errs <- err
					return
				}
				if _, err := redisIsMember(client, setName, members[0]); err != nil {
					errs <- err
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
