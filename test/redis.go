package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

func main() {
	// Connect to redis.
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})
	defer client.Close()

	// Create a new lock client.
	locker := redislock.New(client)

	ctx := context.Background()

	// Try to obtain lock.
	lock, err := locker.Obtain(ctx, "my-key", 30*time.Second, nil)
	if err == redislock.ErrNotObtained {
		fmt.Println("Could not obtain lock!")
		return
	} else if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			ttl, _ := lock.TTL(ctx)
			if ttl < 10*time.Second {
				// Extend my lock.
				if err := lock.Refresh(ctx, 30*time.Second, nil); err != nil {
					log.Fatalln(err)
				}
			}
		}
	}()

	// Don't forget to defer Release.
	defer lock.Release(ctx)
	fmt.Println("I have a lock!")

	// Sleep and check the remaining TTL.
	time.Sleep(20 * time.Second)
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl > 0 {
		fmt.Println("Yay, I still have my lock!")
	}

	// Sleep a little longer, then check.
	time.Sleep(100 * time.Second)
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl == 0 {
		fmt.Println("Now, my lock has expired!")
	}

}
