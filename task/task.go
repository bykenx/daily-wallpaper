package task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
)

var c *cron.Cron

type Task struct {
	Id       cron.EntryID
	Callback func()
	Time     string
}

func (t *Task) Start() {
	if t.Time == "" {
		return
	}
	if t.Id == 0 {
		t.Id, _ = c.AddFunc(getCronSpec(t.Time), t.Callback)
	}
}

func (t *Task) Stop() {
	if t.Id != 0 {
		c.Remove(t.Id)
	}
	t.Id = 0
}

func (t *Task) Restart() {
	t.Stop()
	t.Start()
}

func (t *Task) RunAt(time string) *Task {
	t.Time = time
	return t
}

func NewTask(callback func()) *Task {
	return &Task{
		Callback: callback,
		Time:     "",
	}
}

func init() {
	c = cron.New()
	c.Start()
}

func getCronSpec(time string) string {
	l := strings.Split(time, ":")
	hour := l[0]
	minute := l[1]
	return fmt.Sprintf("%s %s * * *", minute, hour)
}
