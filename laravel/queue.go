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

func (j *QueueJob) OnQueue(q string) *QueueJob {
	j.Queue = q

	return j
}

func (j *QueueJob) WithMaxTries(t int) *QueueJob {
	*j.MaxTries = t

	return j
}

func (j *QueueJob) WithTimeout(t int) *QueueJob {
	*j.Timeout = t

	return j
}

func (j QueueJob) CreateJobPayload() JobPayload {
	data, _ := php.Serialize(j.Payload)

	jobClassString := "O:" + strconv.Itoa(len(j.JobClass)) + ":\"" + j.JobClass + "\":"
	data = jobClassString + data[len("O:3:\"Job\":"):]

	id := uuid.New()

	return JobPayload{
		Uuid:          id.String(),
		DisplayName:   j.JobClass,
		Job:           "Illuminate\\Queue\\CallQueuedHandler@call",
		MaxTries:      j.MaxTries,
		MaxExceptions: nil,
		FailOnTimeout: false,
		Backoff:       false,
		Timeout:       j.Timeout,
		RetryUntil:    nil,
		Data: JobPayloadData{
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
	payload, _ := json.Marshal(job.CreateJobPayload())

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

type JobPayload struct {
	Uuid          string         `json:"uuid"`
	DisplayName   string         `json:"displayName"`
	Job           string         `json:"job"`
	MaxTries      *int           `json:"maxTries"`
	MaxExceptions *int           `json:"maxExceptions"`
	FailOnTimeout bool           `json:"failOnTimeout"`
	Backoff       bool           `json:"backoff"`
	Timeout       *int           `json:"timeout"`
	RetryUntil    *int           `json:"retryUntil"`
	Data          JobPayloadData `json:"data"`
	Id            string         `json:"id"`
	Attepmts      int            `json:"attempts"`
	Type          string         `json:"type"`
	Tags          []string       `json:"tags"`
	Silenced      bool           `json:"slicenced"`
	PushedAt      string         `json:"pushedAt"`
}

type JobPayloadData struct {
	CommandName string `json:"commandName"`
	Command     string `json:"command"`
}
