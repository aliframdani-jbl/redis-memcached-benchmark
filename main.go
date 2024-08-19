package main

import (
	"context"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var memcachedClient *memcache.Client
var ctx = context.Background()

func init() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	// Initialize Memcached client
	memcachedClient = memcache.New("localhost:11211") // Memcached server address
}

// Write and Get Data in Redis
func redisWrite(key, value string) error {
	return redisClient.Set(ctx, key, value, 5*time.Minute).Err()
}

func redisGet(key string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// Write and Get Data in Memcached
func memcachedWrite(key, value string) error {
	return memcachedClient.Set(&memcache.Item{Key: key, Value: []byte(value)})
}

func memcachedGet(key string) (string, error) {
	item, err := memcachedClient.Get(key)
	if err == memcache.ErrCacheMiss {
		return "", nil
	}
	return string(item.Value), err
}

func main() {
	// Example usage
	err := redisWrite("testKey", "testValue")
	if err != nil {
		log.Fatalf("Redis write error: %v", err)
	}
	val, err := redisGet("testKey")
	if err != nil {
		log.Fatalf("Redis get error: %v", err)
	}
	log.Printf("Redis value: %s", val)

	err = memcachedWrite("testKey", "testValue")
	if err != nil {
		log.Fatalf("Memcached write error: %v", err)
	}
	val, err = memcachedGet("testKey")
	if err != nil {
		log.Fatalf("Memcached get error: %v", err)
	}
	log.Printf("Memcached value: %s", val)
}
