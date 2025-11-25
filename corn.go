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

// RunJob 用来在调度器中注册一个任务。
func (corn *CornConfig) RunJob(df gocron.JobDefinition, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(df, task, options...)
}

// redisKey 用来构建任务锁的 Redis 键名。
func (corn *CornConfig) redisKey(id any) string {
	return fmt.Sprintf("corn-%v-lock", id)
}

// RunJobTheOne 用来在持有 Redis 锁时才注册任务，避免重复调度。
func (corn *CornConfig) RunJobTheOne(id any, df gocron.JobDefinition, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(df, task, options...)
	}
	return nil, nil
}

// RunJobEveryDuration 用来按固定间隔周期执行任务。
func (corn *CornConfig) RunJobEveryDuration(duration time.Duration, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.DurationJob(duration),
		task,
		options...,
	)
}

// RunJobEveryDurationTheOne 用来在加锁后以固定间隔执行任务。
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

// RunJobiATime 用来在指定时间运行一次任务。
func (corn *CornConfig) RunJobiATime(time time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), task, options...)
}

// RunJobiATimeTheOne 用来加锁后在指定时间运行一次任务。
func (corn *CornConfig) RunJobiATimeTheOne(id any, time time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time)), task, options...)
	}
	return nil, nil

}

// RunJobiATimes 用来在多个指定时间各运行一次任务。
func (corn *CornConfig) RunJobiATimes(times []time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)), task, options...)
}

// RunJobiATimesTheOne 用来在持锁情况下在多个时间执行任务。
func (corn *CornConfig) RunJobiATimesTheOne(id any, times []time.Time, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	if InsRedis == nil {
		return nil, redisClientNilErr()
	}
	if err := InsRedis.NewLock(corn.redisKey(id)).TryLock(); err == nil {
		return corn.Scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(times...)), task, options...)
	}
	return nil, nil
}

// RunJobEverDay 用来按每天固定时间段执行任务。
func (corn *CornConfig) RunJobEverDay(hours, minutes, seconds, interval uint, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.DailyJob(interval, gocron.NewAtTimes(
			gocron.NewAtTime(hours, minutes, seconds),
		)),
		task,
		options...,
	)
}

// RunJobEverDayTheOne 用来在加锁后每天固定时刻执行任务。
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

// RunJobCrontab 用来根据 cron 表达式调度任务。
func (corn *CornConfig) RunJobCrontab(crontab string, withSeconds bool, task gocron.Task, options ...gocron.JobOption) (gocron.Job, error) {
	return corn.Scheduler.NewJob(
		gocron.CronJob(crontab, withSeconds),
		task,
		options...,
	)
}

// RunJobCrontabTheOne 用来在锁定后按 cron 表达式调度任务。
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

// WithLocation 用来指定调度器使用的时区。
func WithLocation(loc *time.Location) CornOptionFunc {
	return func(c *CornConfig) {
		c.location = loc
	}
}

// WithBeforeJobRuns 用来设置任务执行前的回调。
func WithBeforeJobRuns(beforeJobRuns func(jobID uuid.UUID, jobName string)) CornOptionFunc {
	return func(c *CornConfig) {
		c.beforeJobRuns = beforeJobRuns
	}
}

// WithAfterJobRuns 用来设置任务执行后的回调。
func WithAfterJobRuns(afterJobRuns func(jobID uuid.UUID, jobName string)) CornOptionFunc {
	return func(c *CornConfig) {
		c.afterJobRuns = afterJobRuns
	}
}

// WithAfterJobRunsWithError 用来设置任务出错时的回调。
func WithAfterJobRunsWithError(afterJobRunsWithError func(jobID uuid.UUID, jobName string, err error)) CornOptionFunc {
	return func(c *CornConfig) {
		c.afterJobRunsWithError = afterJobRunsWithError
	}
}

// WithCornJobs 用来追加自定义的调度器选项。
func WithCornJobs(options ...gocron.SchedulerOption) CornOptionFunc {
	return func(c *CornConfig) {
		c.options = append(c.options, options...)
	}
}

// InitCornJob 用来初始化 gocron 调度器并保存全局实例。
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

// Start 用来启动调度器开始运行任务。
func (corn *CornConfig) Start() {
	corn.Scheduler.Start()
}

// Stop 用来优雅关闭调度器。
func (corn *CornConfig) Stop() error {
	if err := corn.Scheduler.Shutdown(); err != nil {
		return err
	}
	return nil
}
