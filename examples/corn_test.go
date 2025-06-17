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

	job, err := gb.CornJob.Scheduler.NewJob(gocron.DurationJob(1*time.Second), gocron.NewTask(func() {
		println("hello world")
	}), gocron.WithLimitedRuns(1))
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(job.ID())

	gb.CornJob.Start()
	defer gb.CornJob.Stop()
	select {}
}
