package reminder

import (
	"fmt"
	"os/exec"
	"time"
)

// Reminder -
type Reminder interface {
	Start()
	Stop()
}

type reminder struct {
	done   chan struct{}
	config Config
}

// Task -
type Task struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Interval string `json:"interval"`
}

// Config -
type Config struct {
	Reminders []Task `json:"reminders"`
}

func notify(title, message string) error {
	osa, err := exec.LookPath("osascript")
	if err != nil {
		return err
	}

	cmd := exec.Command(osa, "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, message, title))

	return cmd.Run()
}

func remind(r <-chan Task, t chan<- Task, d <-chan struct{}) {
	for {
		select {
		case re := <-r:
			notify(re.Title, re.Message)
			t <- re
		case <-d:
			return
		}
	}
}

func ticker(re Task, r chan<- Task, d <-chan struct{}) {
	dur, err := time.ParseDuration(re.Interval)
	if err != nil {
		return
	}

	t := time.Tick(dur)
	for {
		select {
		case <-t:
			r <- re
			return
		case <-d:
			return
		}
	}
}

func timer(t <-chan Task, r chan<- Task, d <-chan struct{}) {
	for {
		select {
		case re := <-t:
			go ticker(re, r, d)
		case <-d:
			return
		}
	}
}

// Start -
func (r *reminder) Start() {
	reminderChan := make(chan Task, len(r.config.Reminders))
	timerChan := make(chan Task, len(r.config.Reminders))

	go remind(reminderChan, timerChan, r.done)

	for _, re := range r.config.Reminders {
		timerChan <- re
	}

	notify("Starting", "Have a good day.")

	timer(timerChan, reminderChan, r.done)
}

// Stop -
func (r *reminder) Stop() {
	close(r.done)
}

// New -
func New(config Config) Reminder {
	return &reminder{
		config: config,
		done:   make(chan struct{}),
	}
}
