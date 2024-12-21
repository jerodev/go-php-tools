package laravel

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jerodev/go-php-tools/php"
	"github.com/redis/go-redis/v9"
)

const broadcastJobClass = `Illuminate\Broadcasting\BroadcastEvent`

type QueueJob struct {
	JobClass string
	MaxTries int
	Payload  interface{}
	Queue    string
	Timeout  int
}

type BroadCastPayload struct {
	Event interface{} `php:"event"`
}

// OnQueue sets the queue name where the job will be dispatched
func (j *QueueJob) OnQueue(q string) *QueueJob {
	j.Queue = q

	return j
}

// WithMaxTries sets the maximum amount of thimes the job will be retried
func (j *QueueJob) WithMaxTries(t int) *QueueJob {
	j.MaxTries = t

	return j
}

// WithTimeout sets the time the job is allowed to run for
func (j *QueueJob) WithTimeout(t int) *QueueJob {
	j.Timeout = t

	return j
}

func (j *QueueJob) createJobPayload() jobPayload {
	php.WithStructNames(map[string]string{
		reflect.ValueOf(j.Payload).Type().Name(): j.JobClass,
	})
	data, _ := php.Serialize(j.Payload)

	id := uuid.New().String()

	payload := jobPayload{
		Uuid:          id,
		DisplayName:   j.JobClass,
		Job:           "Illuminate\\Queue\\CallQueuedHandler@call",
		MaxTries:      nil,
		MaxExceptions: nil,
		FailOnTimeout: false,
		Backoff:       false,
		Timeout:       nil,
		RetryUntil:    nil,
		Data: jobPayloadData{
			CommandName: j.JobClass,
			Command:     data,
		},
		Id:       id,
		Attepmts: 0,
		Type:     "job",
		PushedAt: strconv.Itoa(int(time.Now().UnixMicro())),
	}

	if j.MaxTries > 0 {
		payload.MaxTries = &j.MaxTries
	}

	if j.Timeout > 0 {
		payload.Timeout = &j.Timeout
	}

	return payload
}

func NewBroadcastEvent(jobClass string, payload interface{}) (QueueJob, error) {
	if reflect.ValueOf(payload).Kind() != reflect.Struct {
		return QueueJob{}, errors.New("payload should be a struct")
	}

	php.WithStructNames(map[string]string{
		reflect.ValueOf(payload).Type().Name(): jobClass,
	})

	return NewQueueJob(broadcastJobClass, BroadCastPayload{
		Event: payload,
	})
}

func NewQueueJob(jobClass string, payload interface{}) (QueueJob, error) {
	if reflect.ValueOf(payload).Kind() != reflect.Struct {
		return QueueJob{}, errors.New("payload should be a struct")
	}

	return QueueJob{
		JobClass: jobClass,
		Payload:  payload,
		Queue:    "default",
		Timeout:  0,
	}, nil
}

type QueueConnection interface {
	Dispatch(job QueueJob)
}

type RedisQueueConnection struct {
	client  redis.Client
	context context.Context
	prefix  string
}

func (c *RedisQueueConnection) WithContext(ctx context.Context) {
	c.context = ctx
}

func (c *RedisQueueConnection) Dispatch(job QueueJob) error {
	queueName := c.prefix + job.Queue
	payload, _ := json.Marshal(job.createJobPayload())

	cmd := c.client.RPush(
		c.context,
		queueName,
		string(payload),
	)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	cmd = c.client.RPush(
		c.context,
		queueName+":notify",
		1,
	)

	return cmd.Err()
}

func NewRedisQueueClient(laravelAppName string, opts *redis.Options) RedisQueueConnection {
	return RedisQueueConnection{
		client:  *redis.NewClient(opts),
		context: context.Background(),
		prefix:  strings.ToLower(strings.ReplaceAll(laravelAppName, " ", "_")) + "_database_queues:",
	}
}

type jobPayload struct {
	Uuid          string         `json:"uuid"`
	DisplayName   string         `json:"displayName"`
	Job           string         `json:"job"`
	MaxTries      *int           `json:"maxTries"`
	MaxExceptions *int           `json:"maxExceptions"`
	FailOnTimeout bool           `json:"failOnTimeout"`
	Backoff       bool           `json:"backoff"`
	Timeout       *int           `json:"timeout"`
	RetryUntil    *int           `json:"retryUntil"`
	Data          jobPayloadData `json:"data"`
	Id            string         `json:"id"`
	Attepmts      int            `json:"attempts"`
	Type          string         `json:"type"`
	Tags          []string       `json:"tags"`
	Silenced      bool           `json:"slicenced"`
	PushedAt      string         `json:"pushedAt"`
}

type jobPayloadData struct {
	CommandName string `json:"commandName"`
	Command     string `json:"command"`
}
