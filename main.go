package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func acquireLock(client *redis.Client, key string, expiration time.Duration) (bool, error) {
	// Use the SET command to try to acquire the lock
	result, err := client.SetNX(key, "lock", expiration).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func releaseLock(client *redis.Client, key string) error {
	// Use the DEL command to release the lock
	_, err := client.Del(key).Result()
	return err
}

func main() {
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	mutexKey := "my_lock"

	// Try to acquire the first lock
	isFirstLockSet, err := acquireLock(client, mutexKey, time.Second*10)
	if err != nil {
		fmt.Println("Error acquiring lock:", err)
		return
	}
	if !isFirstLockSet {
		fmt.Println("Failed to acquire lock")
		return
	}

	// Do some work while holding the lock
	fmt.Println("First lock acquired!")

	isSecondLockSet, _ := acquireLock(client, mutexKey, time.Second*10)

	if !isSecondLockSet {
		fmt.Println("Could not get a second lock which is as expected this is where you would force the request out.")
	} else {
		fmt.Println("Second Lock acquired! This should not happen :)")
	}

	// Simulate some work by sleeping and try to acquire the lock again to see that it fails
	time.Sleep(time.Second * 5)

	isThirdLockSet, _ := acquireLock(client, mutexKey, time.Second*10)
	if !isThirdLockSet {
		fmt.Println("Still could not get the third lock since the first lock is still set.")
	} else {
		fmt.Println("Third Lock acquired! This should not happen :)")
	}

	// Release the lock
	err = releaseLock(client, mutexKey)
	if err != nil {
		fmt.Println("Error releasing lock:", err)
		return
	}
	fmt.Println("First lock released!")

	// Try to acquire the lock again to show it has been released
	isForthLockSet, _ := acquireLock(client, mutexKey, time.Second*10)
	if !isForthLockSet {
		fmt.Println("Failed to acquire lock")
		return
	} else {
		fmt.Println("Forth Lock acquired!")
	}

	// Release the forth lock
	err = releaseLock(client, mutexKey)
	if err != nil {
		fmt.Println("Error releasing lock:", err)
		return
	}
	fmt.Println("Forth Lock released!")
}
