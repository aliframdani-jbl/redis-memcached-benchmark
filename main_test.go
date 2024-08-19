package main

import (
	"sync"
	"testing"
)

// Benchmark testing with concurrency
func BenchmarkRedis(b *testing.B) {
	const numGoroutines = 100
	const numIterations = 100000

	var wg sync.WaitGroup
	errs := make(chan error, numGoroutines) // Channel to collect errors
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations/numGoroutines; j++ {
				if err := redisWrite("key", "value"); err != nil {
					errs <- err
					return
				}
				if _, err := redisGet("key"); err != nil {
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
		b.Fatal(err)
	}
}

func BenchmarkMemcached(b *testing.B) {
	const numGoroutines = 100
	const numIterations = 100000

	var wg sync.WaitGroup
	errs := make(chan error, numGoroutines) // Channel to collect errors
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations/numGoroutines; j++ {
				if err := memcachedWrite("key", "value"); err != nil {
					errs <- err
					return
				}
				if _, err := memcachedGet("key"); err != nil {
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
		b.Fatal(err)
	}
}
