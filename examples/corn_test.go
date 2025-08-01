package examples

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/loveyu233/gb"
	"testing"
	"time"
)

func TestCorn(t *testing.T) {
	err := gb.InitCornJob()
	if err != nil {
		t.Fatal(err)
		return
	}

	//job, err := gb.CornJob.Scheduler.NewJob(gocron.DurationJob(1*time.Second), gocron.NewTask(func() {
	//	println("hello world")
	//}), gocron.WithLimitedRuns(1))
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}
	//t.Log(job.ID())

	//gb.CornJob.RunJobiATime(time.Now().Add(3*time.Second), gocron.NewTask(func(s string, i int) {
	//	t.Log(time.Now(), s, " ", i)
	//}, "1", 1))

	gb.CornJob.RunJobiATimes([]time.Time{time.Now().Add(3 * time.Second), time.Now().Add(5 * time.Second)}, gocron.NewTask(func(s string, i int) {
		t.Log(time.Now(), s, " ", i)
	}, "1", 1))

	gb.CornJob.Start()
	defer gb.CornJob.Stop()
	select {}
}
