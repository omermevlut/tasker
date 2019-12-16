package tasker

import "time"

var day = 24 * time.Hour
var week = 7 * day

// manages schedule times of the tasks
type timer struct {
	IsRepeating    bool         `json:"is_recurring"`
	Hour           int64        `json:"hour"`
	Minute         int64        `json:"minute"`
	WeekDay        time.Weekday `json:"week_day"`
	OccurrenceType string       `json:"occurrence_type"`
	RunAt          int64        `json:"run_at"`
	IsInfinite     bool         `json:"is_infinite"`
	UntilTime      time.Time    `json:"until_time"`
	RunCount       int64        `json:"run_count"`
	MaxRunCount    int64        `json:"max_run_count"`
}

func (t *timer) createNextDailyRunDate() {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), int(t.Hour), int(t.Minute), 0, 0, time.Local)

	if tm.Unix() <= now.Unix() {
		tm = tm.Add(day)
	}

	t.OccurrenceType = "daily"
	t.RunAt = tm.Unix()
}

func (t *timer) createNextHourlyRunDate() {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), int(t.Minute), 0, 0, time.Local)

	if tm.Unix() <= now.Unix() {
		tm = tm.Add(time.Hour)
	}

	t.OccurrenceType = "hourly"
	t.RunAt = tm.Unix()
}

func (t *timer) createNextWeeklyRunDate() {
	now := time.Now()

	for {
		if int(now.Weekday()) == int(t.WeekDay) {
			break
		}

		now = now.Add(day)
	}

	tm := time.Date(now.Year(), now.Month(), now.Day(), int(t.Hour), int(t.Minute), 0, 0, time.Local)

	if tm.Unix() <= time.Now().Unix() {
		tm = tm.Add(week)
	}

	t.OccurrenceType = "weekly"
	t.RunAt = tm.Unix()
}

func (t *timer) setNextRun() {
	switch t.OccurrenceType {
	case "weekly":
		t.createNextWeeklyRunDate()
	case "daily":
		t.createNextDailyRunDate()
	case "hourly":
		t.createNextHourlyRunDate()
	}
}
