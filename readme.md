# Go PHP Tools
A bundle of tools for interaction between Go and PHP written in Go

- [PHP Tools](#php-tools)
- [Laravel Tools](#laravel-tools)
  - [Queue Redis job](#queue-redis-job)
  - [Queue Broadcasting event](#queue-broadcasting-event)

## PHP Tools
### Serialize data
The `serialize()` functions takes a go variable and serializes it to PHP serialization format.

```go
package main

import (
	"fmt"
	"github.com/jerodev/go-php-tools/php"
)

type User struct {
	Name string `php:"username"`
	Age  int    `php:"age"`
}

func main() {
	php, _ := php.Serialize(User{
		Name: "Jerodev",
		Age:  30,
	})

	fmt.Println(php) // O:4:"User":2:{s:8:"username";s:7:"Jerodev";s:3:"age";i:30;}
}
```

## Laravel Tools
### Queue Redis job
Queue a job on a queue that can be executed by a laravel queue worker.

The struct passed to the `NewQueueJob` function should match the data required in the actual PHP job class.

```go
package main

import (
	"github.com/jerodev/go-php-tools/laravel"
	"github.com/redis/go-redis/v9"
)

type Job struct {
	Contents string `php:"contents"`
}

func main() {
	job, _ := laravel.NewQueueJob("App\\Jobs\\LaravelTestJob", Job{
		Contents: "Lorem Ipsum",
	})

	conn := laravel.NewRedisQueueClient("LaraQueue", &redis.Options{
		Addr: "127.0.0.1:6379",
	})

	conn.Dispatch(job)
}

```

### Queue Broadcasting event
In the same way that this package can queue jobs, it can also queue broadcasting events to be picked up by your Laravel
application or Laravel Reverb.

```go
package main

import (
	"github.com/jerodev/go-php-tools/laravel"
	"github.com/redis/go-redis/v9"
)

type OrderUpdated struct {
	OrderId int `php:"order_id"`
}

func main() {
	job, _ := laravel.NewBroadcastEvent("App\\Events\\OrderUpdated", OrderUpdated{
		OrderId: 42,
	})

	conn := laravel.NewRedisQueueClient("LaraQueue", &redis.Options{
		Addr: "127.0.0.1:6379",
	})

	conn.Dispatch(job)
}
```