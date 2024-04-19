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

type QueueJob struct {
	JobClass string
	MaxTries *int
	Payload  interface{}
	Queue    string
	Timeout  *int
}

// OnQueue sets the queue name where the job will be dispatched
func (j *QueueJob) OnQueue(q string) *QueueJob {
	j.Queue = q

	return j
}

// WithMaxTries sets the maximum amount of thimes the job will be retried
func (j *QueueJob) WithMaxTries(t int) *QueueJob {
	*j.MaxTries = t

	return j
}

// WithTimeout sets the time the job is allowed to run for
func (j *QueueJob) WithTimeout(t int) *QueueJob {
	*j.Timeout = t

	return j
}

func (j QueueJob) createJobPayload() jobPayload {
	refl := reflect.ValueOf(j.Payload)
	data, _ := php.Serialize(j.Payload)

	jobClassString := "O:" + strconv.Itoa(len(j.JobClass)) + ":\"" + j.JobClass + "\":"
	data = jobClassString + data[len(refl.Type().Name())+len("O:3:\"\":"):]

	id := uuid.New().String()

	return jobPayload{
		Uuid:          id,
		DisplayName:   j.JobClass,
		Job:           "Illuminate\\Queue\\CallQueuedHandler@call",
		MaxTries:      j.MaxTries,
		MaxExceptions: nil,
		FailOnTimeout: false,
		Backoff:       false,
		Timeout:       j.Timeout,
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
}

func NewQueueJob(jobClass string, payload interface{}) (*QueueJob, error) {
	if reflect.ValueOf(payload).Kind() != reflect.Struct {
		return nil, errors.New("payload should be a struct")
	}

	return &QueueJob{
		JobClass: jobClass,
		Payload:  payload,
		Queue:    "default",
		Timeout:  nil,
	}, nil
}

type QueueConnection interface {
	Dispatch(job QueueJob)
}

type RedisQueueConnection struct {
	client redis.Client
	prefix string
}

func (c RedisQueueConnection) Dispatch(job QueueJob) {
	queueName := c.prefix + job.Queue
	payload, _ := json.Marshal(job.createJobPayload())

	ctx := context.Background()

	c.client.RPush(
		ctx,
		queueName,
		string(payload),
	)

	c.client.RPush(
		ctx,
		queueName+":notify",
		1,
	)
}

func NewRedisQueueClient(laravelAppName string, opts *redis.Options) RedisQueueConnection {
	return RedisQueueConnection{
		client: *redis.NewClient(opts),
		prefix: strings.ToLower(strings.ReplaceAll(laravelAppName, " ", "_")) + "_database_queues:",
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
