package main

import (
	"os"

	"github.com/jerodev/go-php-tools/laravel"
	"github.com/redis/go-redis/v9"
)

type Job struct {
	Contents string `php:"contents"`
}

func main() {
	job, _ := laravel.NewQueueJob("App\\Jobs\\LaravelTestJob", Job{
		Contents: os.Getenv("TEST_CONTENT"),
	})

	conn := laravel.NewRedisQueueClient("LaraQueue", &redis.Options{
		Addr: "127.0.0.1:6379",
	})

	conn.Dispatch(job)
}
