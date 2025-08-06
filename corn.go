package gb

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"time"
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

// RunJob 运行自定义的定时任务
func (corn *CornConfig) RunJob(df gocron.JobDefinition, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(df, task, options...)
}
func (corn *CornConfig) redisKey(id any) string {
	return fmt.Sprintf("corn-%v-lock", id)
}

// RunJobTheOne 运行自定义的定时任务,使用redis分布式锁进行控制同一id只会运行一个任务
func (corn *CornConfig) RunJobTheOne(id any, df gocron.JobDefinition, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(df, task, options...)
	}
	return nil, nil
}

// RunJobEveryDuration 创建每duration时间执行一次的定时任务
func (corn *CornConfig) RunJobEveryDuration(duration time.Duration, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.DurationJob(duration),
		task,
		options...,
	)
}

// RunJobEveryDurationTheOne 使用redis分布式锁进行控制同一id只会运行一个任务
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

// RunJobiATime 运行指定时间的定时任务
func (corn *CornConfig) RunJobiATime(time time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), task, options...)
}

// RunJobiATimeTheOne 使用redis分布式锁进行控制同一id只会运行一个任务
func (corn *CornConfig) RunJobiATimeTheOne(id any, time time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), task, options...)
	}
	return nil, nil

}

// RunJobiATimes 运行多个时间的定时任务
func (corn *CornConfig) RunJobiATimes(times []time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)), task, options...)
}

// RunJobiATimesTheOne 使用redis分布式锁进行控制同一id只会运行一个任务
func (corn *CornConfig) RunJobiATimesTheOne(id any, times []time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)), task, options...)
	}
	return nil, nil
}

// RunJobEverDay 在指定的时间内每天运行interval次task
func (corn *CornConfig) RunJobEverDay(hours, minutes, seconds, interval uint, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.DailyJob(interval, gocron.NewAtTimes(
			gocron.NewAtTime(hours, minutes, seconds),
		)),
		task,
		options...,
	)
}

// RunJobEverDayTheOne 使用redis分布式锁进行控制同一id只会运行一个任务
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

// RunJobCrontab 使用 Cron 表达式,如果withSeconds设置为true，则可以在开始时使用可选的第6个字段
func (corn *CornConfig) RunJobCrontab(crontab string, withSeconds bool, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.CronJob(crontab, withSeconds),
		task,
		options...,
	)
}

// RunJobCrontabTheOne 使用redis分布式锁进行控制同一id只会运行一个任务
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

func WithLocation(loc *time.Location) CornOptionFunc {
	return func(c *CornConfig) {
		c.location = loc
	}
}

func WithBeforeJobRuns(beforeJobRuns func(jobID uuid.UUID, jobName string)) CornOptionFunc {
	return func(c *CornConfig) {
		c.beforeJobRuns = beforeJobRuns
	}
}

func WithAfterJobRuns(afterJobRuns func(jobID uuid.UUID, jobName string)) CornOptionFunc {
	return func(c *CornConfig) {
		c.afterJobRuns = afterJobRuns
	}
}

func WithAfterJobRunsWithError(afterJobRunsWithError func(jobID uuid.UUID, jobName string, err error)) CornOptionFunc {
	return func(c *CornConfig) {
		c.afterJobRunsWithError = afterJobRunsWithError
	}
}

func WithCornJobs(options ...gocron.SchedulerOption) CornOptionFunc {
	return func(c *CornConfig) {
		c.options = append(c.options, options...)
	}
}

// InitCornJob 初始化定时任务,默认时区使用ShangHaiTimeLocation
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

// Start 开启定时任务
func (corn *CornConfig) Start() {
	corn.Scheduler.Start()
}

// Stop 结束定时任务
func (corn *CornConfig) Stop() error {
	if err := corn.Scheduler.Shutdown(); err != nil {
		return err
	}
	return nil
}
