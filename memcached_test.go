package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

func memcachedSet(client *memcache.Client, key, value string) error {
	return client.Set(&memcache.Item{Key: key, Value: []byte(value)})
}

func memcachedGet(client *memcache.Client, key string) (string, error) {
	item, err := client.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

func BenchmarkMemcachedSetGet(b *testing.B) {
	client := memcache.New("localhost:11211")
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
				if err := memcachedSet(client, key, value); err != nil {
					errs <- err
					return
				}
				retrievedValue, err := memcachedGet(client, key)
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

func randMarketplaceAndStoreId() string {
	var id string
	for i := 0; i < 8; i++ {
		id = id + fmt.Sprintf("%d", rand.Intn(10)) // Generates a random integer between 0 and 9
	}

	return id
}

func randMarketplace() string {
	marketplaces := []string{
		"SHOPEE",
		"LAZADA",
		"TOKOPEDIA",
		"BLIBLI",
		"BUKALAPAK",
		"TIKTOK",
		"SHOPIFY",
		"WOOCOMMERCE",
		"ZALORA",
	}

	randomIndex := rand.Intn(len(marketplaces))
	return marketplaces[randomIndex]
}
