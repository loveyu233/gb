package gb

import (
	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var CornJob *Corn

type Corn struct {
	location              *time.Location                                   // 时区
	beforeJobRuns         func(jobID uuid.UUID, jobName string)            // 运行前
	afterJobRuns          func(jobID uuid.UUID, jobName string)            // 运行后
	afterJobRunsWithError func(jobID uuid.UUID, jobName string, err error) // 出错
	options               []gocron.SchedulerOption
	redisClient           *redis.Client
	Scheduler             gocron.Scheduler
}

type CornOption func(*Corn)

func WithLocation(loc *time.Location) CornOption {
	return func(c *Corn) {
		c.location = loc
	}
}

func WithBeforeJobRuns(beforeJobRuns func(jobID uuid.UUID, jobName string)) CornOption {
	return func(c *Corn) {
		c.beforeJobRuns = beforeJobRuns
	}
}

func WithAfterJobRuns(afterJobRuns func(jobID uuid.UUID, jobName string)) CornOption {
	return func(c *Corn) {
		c.afterJobRuns = afterJobRuns
	}
}

func WithAfterJobRunsWithError(afterJobRunsWithError func(jobID uuid.UUID, jobName string, err error)) CornOption {
	return func(c *Corn) {
		c.afterJobRunsWithError = afterJobRunsWithError
	}
}

func WithCornJobs(options ...gocron.SchedulerOption) CornOption {
	return func(c *Corn) {
		c.options = append(c.options, options...)
	}
}

func WithCornRedisClient(client *redis.Client) CornOption {
	return func(c *Corn) {
		c.redisClient = client
	}
}

func NewCornJob(options ...CornOption) error {
	var corn = &Corn{
		options: make([]gocron.SchedulerOption, 0),
	}
	for _, opt := range options {
		opt(corn)
	}

	if corn.location == nil {
		cst, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			return err
		}
		corn.location = cst
	}
	corn.options = append(corn.options, gocron.WithLocation(corn.location))

	var eventListeners []gocron.EventListener
	if corn.afterJobRuns != nil {
		eventListeners = append(eventListeners, gocron.AfterJobRuns(corn.afterJobRuns))
	}
	if corn.beforeJobRuns != nil {
		eventListeners = append(eventListeners, gocron.BeforeJobRuns(corn.beforeJobRuns))
	}
	if corn.afterJobRunsWithError != nil {
		eventListeners = append(eventListeners, gocron.AfterJobRunsWithError(corn.afterJobRunsWithError))
	}
	if len(eventListeners) > 0 {
		corn.options = append(corn.options, gocron.WithGlobalJobOptions(
			gocron.WithEventListeners(eventListeners...),
		))
	}

	if corn.redisClient != nil {
		locker, err := redislock.NewRedisLocker(corn.redisClient)
		if err != nil {
			return err
		}
		corn.options = append(corn.options, gocron.WithDistributedLocker(locker))
	}

	scheduler, err := gocron.NewScheduler(corn.options...)
	if err != nil {
		return err
	}

	corn.Scheduler = scheduler
	CornJob = corn
	return nil
}

func (corn *Corn) Start() {
	corn.Scheduler.Start()
}

func (corn *Corn) Stop() error {
	if err := corn.Scheduler.Shutdown(); err != nil {
		return err
	}
	return nil
}
