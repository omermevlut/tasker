package tasker

import (
	"time"

	"github.com/google/uuid"
)

// TaskData ...
type TaskData map[string]interface{}

// Task ...
type Task struct {
	timer
	validator

	Data       *TaskData `json:"task_data"`
	Name       string    `json:"name"`
	ID         string    `json:"id"`
	ExecutedAt time.Time `json:"executed_at"`
}

// NewTask ...
func NewTask(name string, data *TaskData) *Task {
	return &Task{
		Data: data,
		Name: name,
		ID:   uuid.New().String(),
	}
}

// OnceAt task
func (t *Task) OnceAt(expiry int64) *Task {
	t.RunAt = expiry

	return t
}

// DailyAt runs the task daily at the given hour and minute
//
// hour: `min:"0" max:"23"`,
// minute: `min:"0" max:"59"`
func (t *Task) DailyAt(hour, minute int64) *Task {
	t.validateHour(hour)
	t.validateMinute(minute)

	t.Hour = hour
	t.Minute = minute
	t.IsRepeating = true
	t.IsInfinite = true

	t.createNextDailyRunDate()

	return t
}

// HourlyAt runs the task hourly at the given minute
//
// minute: `min:"0" max:"59"`
func (t *Task) HourlyAt(minute int64) *Task {
	t.validateMinute(minute)

	t.Minute = minute
	t.IsRepeating = true
	t.IsInfinite = true

	t.createNextHourlyRunDate()

	return t
}

// WeeklyAt runs the task periodically every week
func (t *Task) WeeklyAt(weekday time.Weekday, hour, minute int64) *Task {
	t.validateHour(hour)
	t.validateMinute(minute)

	t.Minute = minute
	t.Hour = hour
	t.WeekDay = weekday
	t.IsRepeating = true
	t.IsInfinite = true

	t.createNextWeeklyRunDate()

	return t
}

// Until repeats the task until the given time
func (t *Task) Until(until time.Time) *Task {
	t.UntilTime = until
	t.IsInfinite = false

	return t
}

// UntilCount repeats until count `r` is reached
func (t *Task) UntilCount(r int64) *Task {
	t.MaxRunCount = r
	t.IsInfinite = false

	return t
}

func (t *Task) isExpired() bool {
	return (t.UntilTime.Unix() < time.Now().Unix()) && !t.IsInfinite
}
