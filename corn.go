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

// RunJob 运行自定义的定时任务
func (corn *Corn) RunJob(df gocron.JobDefinition, task gocron.Task, opts ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(df, task, opts...)
}

// RunJobiATime 运行指定时间的定时任务
func (corn *Corn) RunJobiATime(time time.Time, task gocron.Task, parameters ...any) (gocron.Job, error) {
	return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), task)
}

// RunJobiATimes 运行多个时间的定时任务
func (corn *Corn) RunJobiATimes(times []time.Time, task gocron.Task, parameters ...any) (gocron.Job, error) {
	return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)), task)
}

// RunJobEverDay 在指定的时间内每天运行interval次task
func (corn *Corn) RunJobEverDay(hours, minutes, seconds, interval uint, task gocron.Task) (gocron.Job, error) {
	job, err := corn.Scheduler.NewJob(
		gocron.DailyJob(interval, gocron.NewAtTimes(
			gocron.NewAtTime(hours, minutes, seconds), // 5点0分0秒
		)),
		task,
	)
	return job, err
}

// RunJobCrontab 使用 Cron 表达式,如果withSeconds设置为true，则可以在开始时使用可选的第6个字段
func (corn *Corn) RunJobCrontab(crontab string, withSeconds bool, task gocron.Task) (gocron.Job, error) {
	job, err := corn.Scheduler.NewJob(
		gocron.CronJob(crontab, withSeconds), // 每天5点执行
		task,
	)
	return job, err
}

type CornOptionFunc func(*Corn)

func WithLocation(loc *time.Location) CornOptionFunc {
	return func(c *Corn) {
		c.location = loc
	}
}

func WithBeforeJobRuns(beforeJobRuns func(jobID uuid.UUID, jobName string)) CornOptionFunc {
	return func(c *Corn) {
		c.beforeJobRuns = beforeJobRuns
	}
}

func WithAfterJobRuns(afterJobRuns func(jobID uuid.UUID, jobName string)) CornOptionFunc {
	return func(c *Corn) {
		c.afterJobRuns = afterJobRuns
	}
}

func WithAfterJobRunsWithError(afterJobRunsWithError func(jobID uuid.UUID, jobName string, err error)) CornOptionFunc {
	return func(c *Corn) {
		c.afterJobRunsWithError = afterJobRunsWithError
	}
}

func WithCornJobs(options ...gocron.SchedulerOption) CornOptionFunc {
	return func(c *Corn) {
		c.options = append(c.options, options...)
	}
}

func WithCornRedisClient(client *redis.Client) CornOptionFunc {
	return func(c *Corn) {
		c.redisClient = client
	}
}

func InitCornJob(options ...CornOptionFunc) error {
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
