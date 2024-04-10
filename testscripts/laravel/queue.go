package main

import (
	"github.com/jerodev/go-php-tools/laravel"
	"github.com/redis/go-redis/v9"
)

func main() {
	job := laravel.NewQueueJob("App\\Jobs\\LaravelTestJob", nil)

	conn := laravel.NewRedisQueueClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	conn.Dispatch(*job)
}
