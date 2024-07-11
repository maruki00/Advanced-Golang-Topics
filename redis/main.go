package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	err := client.Set(ctx, "key", "value, kldfjhg ,n sdfghskldfjh", 0).Err()
	if err != nil {
		panic("Could not Save The Item,")
	}
	val, err := client.Get(ctx, "key").Result()
	if err != nil {
		panic("Could not get the Item.")
	}
	fmt.Println("Resul: ", val)

}
