package task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strconv"
	"strings"
	"time"
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
	if isArriveAtTime(t.Time) {
		t.Callback()
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

func isArriveAtTime(t string) bool {
	p := strings.Split(t, "")
	n := time.Now()
	h, _ := strconv.Atoi(p[0])
	m, _ := strconv.Atoi(p[1])
	t1 := h * 60 + m
	t2 := n.Hour() * 60 + m
	return t2 >= t1
}