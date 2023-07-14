package tasks

import (
	"fmt"
	"time"
)

type TaskFunc func()

type Task struct {
	fn           TaskFunc
	duration     time.Duration
	hourToRun    int
	minutesToRun int
	secondToRun  int
}

type TaskScheduler struct {
	tasks []Task
}

type TaskTimer struct {
	timer *time.Timer
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{tasks: []Task{}}
}

func (ts *TaskScheduler) AddTask(
	t TaskFunc,
	duration time.Duration,
	hourToRun int,
	minutesToRun int,
	secondToRun int,
) {
	ts.tasks = append(ts.tasks, Task{
		t, duration,
		hourToRun, minutesToRun, secondToRun,
	})
}

func (ts *TaskScheduler) taskRunner(
	fn TaskFunc,
	t *TaskTimer,
	d time.Duration,
	hourToRun int,
	minutesToRun int,
	secondToRun int,
) {
	updateTimerFunc := updateTimer(t, d, hourToRun, minutesToRun, secondToRun)
	for {
		// updateTimer(t, d, hourToRun, minutesToRun, secondToRun)
		updateTimerFunc()
		<-t.timer.C
		fn()
	}
}

func (ts *TaskScheduler) Run() {
	for _, t := range ts.tasks {
		taskTimer := &TaskTimer{timer: time.NewTimer(t.duration)}
		go ts.taskRunner(t.fn, taskTimer, t.duration, t.hourToRun, t.minutesToRun, t.secondToRun)
	}
}

func updateTimer(
	t *TaskTimer,
	d time.Duration,
	hourToRun int,
	minutesToRun int,
	secondToRun int,
) func() {
	nextToRun := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), hourToRun, minutesToRun, secondToRun, 0, time.Local)
	return func() {
		if !nextToRun.After(time.Now()) {
			fmt.Println(nextToRun, "- next tick --------")

			fmt.Println("heeeyyy")
			if !nextToRun.Add(d).After(time.Now()) {
				nextToRun = time.Now().Add(d)
			} else {
				nextToRun = nextToRun.Add(d)
			}
		}
		fmt.Println(nextToRun, "- next tick")
		diff := nextToRun.Sub(time.Now())
		fmt.Println(diff, "[[[[[[[[[[[[[[[[[[]]]]]]]]]]]]]]]]]]")
		t.timer = time.NewTimer(diff)
		t.timer.Reset(diff)
	}

}
