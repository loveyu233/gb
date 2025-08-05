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
	gb.CornJob.RunJobEveryDuration(10*time.Second, gocron.NewTask(func() {
		t.Log(gb.FormatCurrentTime())
	}))
	time.Sleep(time.Second)
	gb.CornJob.Start()
	select {}
}
