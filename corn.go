package gb

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

var InsCornJob *CornConfig

type CornConfig struct {
	location              *time.Location                                   // 时区
	beforeJobRuns         func(jobID uuid.UUID, jobName string)            // 运行前
	afterJobRuns          func(jobID uuid.UUID, jobName string)            // 运行后
	afterJobRunsWithError func(jobID uuid.UUID, jobName string, err error) // 出错
	options               []gocron.SchedulerOption
	Scheduler             gocron.Scheduler
}

// RunJob 方法用于处理RunJob相关逻辑。
func (corn *CornConfig) RunJob(df gocron.JobDefinition, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(df, task, options...)
}

// redisKey 方法用于处理redisKey相关逻辑。
func (corn *CornConfig) redisKey(id any) string {
	return fmt.Sprintf("corn-%v-lock", id)
}

// RunJobTheOne 方法用于处理RunJobTheOne相关逻辑。
func (corn *CornConfig) RunJobTheOne(id any, df gocron.JobDefinition, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(df, task, options...)
	}
	return nil, nil
}

// RunJobEveryDuration 方法用于处理RunJobEveryDuration相关逻辑。
func (corn *CornConfig) RunJobEveryDuration(duration time.Duration, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.DurationJob(duration),
		task,
		options...,
	)
}

// RunJobEveryDurationTheOne 方法用于处理RunJobEveryDurationTheOne相关逻辑。
func (corn *CornConfig) RunJobEveryDurationTheOne(id any, duration time.Duration, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(
			gocron.DurationJob(duration),
			task,
			options...,
		)
	}
	return nil, nil
}

// RunJobiATime 方法用于处理RunJobiATime相关逻辑。
func (corn *CornConfig) RunJobiATime(time time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), task, options...)
}

// RunJobiATimeTheOne 方法用于处理RunJobiATimeTheOne相关逻辑。
func (corn *CornConfig) RunJobiATimeTheOne(id any, time time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), task, options...)
	}
	return nil, nil

}

// RunJobiATimes 方法用于处理RunJobiATimes相关逻辑。
func (corn *CornConfig) RunJobiATimes(times []time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)), task, options...)
}

// RunJobiATimesTheOne 方法用于处理RunJobiATimesTheOne相关逻辑。
func (corn *CornConfig) RunJobiATimesTheOne(id any, times []time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)), task, options...)
	}
	return nil, nil
}

// RunJobEverDay 方法用于处理RunJobEverDay相关逻辑。
func (corn *CornConfig) RunJobEverDay(hours, minutes, seconds, interval uint, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.DailyJob(interval, gocron.NewAtTimes(
			gocron.NewAtTime(hours, minutes, seconds),
		)),
		task,
		options...,
	)
}

// RunJobEverDayTheOne 方法用于处理RunJobEverDayTheOne相关逻辑。
func (corn *CornConfig) RunJobEverDayTheOne(id any, hours, minutes, seconds, interval uint, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(
			gocron.DailyJob(interval, gocron.NewAtTimes(
				gocron.NewAtTime(hours, minutes, seconds),
			)),
			task,
			options...,
		)
	}
	return nil, nil
}

// RunJobCrontab 方法用于处理RunJobCrontab相关逻辑。
func (corn *CornConfig) RunJobCrontab(crontab string, withSeconds bool, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.CronJob(crontab, withSeconds),
		task,
		options...,
	)
}

// RunJobCrontabTheOne 方法用于处理RunJobCrontabTheOne相关逻辑。
func (corn *CornConfig) RunJobCrontabTheOne(id any, crontab string, withSeconds bool, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(
			gocron.CronJob(crontab, withSeconds),
			task,
			options...,
		)
	}
	return nil, nil
}

type CornOptionFunc func(*CornConfig)

// WithLocation 函数用于处理WithLocation相关逻辑。
func WithLocation(loc *time.Location) CornOptionFunc {
	return func(c *CornConfig) {
		c.location = loc
	}
}

// WithBeforeJobRuns 函数用于处理WithBeforeJobRuns相关逻辑。
func WithBeforeJobRuns(beforeJobRuns func(jobID uuid.UUID, jobName string)) CornOptionFunc {
	return func(c *CornConfig) {
		c.beforeJobRuns = beforeJobRuns
	}
}

// WithAfterJobRuns 函数用于处理WithAfterJobRuns相关逻辑。
func WithAfterJobRuns(afterJobRuns func(jobID uuid.UUID, jobName string)) CornOptionFunc {
	return func(c *CornConfig) {
		c.afterJobRuns = afterJobRuns
	}
}

// WithAfterJobRunsWithError 函数用于处理WithAfterJobRunsWithError相关逻辑。
func WithAfterJobRunsWithError(afterJobRunsWithError func(jobID uuid.UUID, jobName string, err error)) CornOptionFunc {
	return func(c *CornConfig) {
		c.afterJobRunsWithError = afterJobRunsWithError
	}
}

// WithCornJobs 函数用于处理WithCornJobs相关逻辑。
func WithCornJobs(options ...gocron.SchedulerOption) CornOptionFunc {
	return func(c *CornConfig) {
		c.options = append(c.options, options...)
	}
}

// InitCornJob 函数用于处理InitCornJob相关逻辑。
func InitCornJob(options ...CornOptionFunc) error {
	var corn = &CornConfig{
		options: make([]gocron.SchedulerOption, 0),
	}
	for _, opt := range options {
		opt(corn)
	}

	if corn.location == nil {
		corn.location = ShangHaiTimeLocation
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

	scheduler, err := gocron.NewScheduler(corn.options...)
	if err != nil {
		return err
	}

	corn.Scheduler = scheduler
	InsCornJob = corn
	return nil
}

// Start 方法用于处理Start相关逻辑。
func (corn *CornConfig) Start() {
	corn.Scheduler.Start()
}

// Stop 方法用于处理Stop相关逻辑。
func (corn *CornConfig) Stop() error {
	if err := corn.Scheduler.Shutdown(); err != nil {
		return err
	}
	return nil
}
