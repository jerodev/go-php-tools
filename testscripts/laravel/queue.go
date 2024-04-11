package main

import (
	"github.com/jerodev/go-php-tools/laravel"
	"github.com/redis/go-redis/v9"
)

type Job struct {
	Contents string `php:"contents"`
}

func main() {
	job := laravel.NewQueueJob("App\\Jobs\\LaravelTestJob", Job{
		Contents: "Lorem Ipsum",
	})

	conn := laravel.NewRedisQueueClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	conn.Dispatch(*job)
}
