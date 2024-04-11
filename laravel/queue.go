package laravel

import (
	"context"
	"encoding/json"
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
	data, _ := php.Serialize(j.Payload)

	jobClassString := "O:" + strconv.Itoa(len(j.JobClass)) + ":\"" + j.JobClass + "\":"
	data = jobClassString + data[len("O:3:\"Job\":"):]

	id := uuid.New()

	return jobPayload{
		Uuid:          id.String(),
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
		Id:       id.String(),
		Attepmts: 0,
		Type:     "job",
		PushedAt: strconv.Itoa(int(time.Now().UnixMicro())),
	}
}

func NewQueueJob(jobClass string, payload interface{}) *QueueJob {
	return &QueueJob{
		JobClass: jobClass,
		Payload:  payload,
		Queue:    "default",
		Timeout:  nil,
	}
}

type QueueConnection interface {
	Dispatch(job QueueJob)
}

type RedisQueueConnection struct {
	client         redis.Client
	laravelAppName string
}

func (c RedisQueueConnection) Dispatch(job QueueJob) {
	queueName := strings.ToLower(strings.ReplaceAll(c.laravelAppName, " ", "_")) + "_database_queues:" + job.Queue
	payload, _ := json.Marshal(job.createJobPayload())

	c.client.RPush(
		context.Background(),
		queueName,
		string(payload),
	)

	c.client.RPush(
		context.Background(),
		queueName+":notify",
		1,
	)
}

func NewRedisQueueClient(laravelAppName string, opts *redis.Options) RedisQueueConnection {
	return RedisQueueConnection{
		client:         *redis.NewClient(opts),
		laravelAppName: laravelAppName,
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
