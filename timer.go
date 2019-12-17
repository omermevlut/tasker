package tasker

import "time"

var day = 24 * time.Hour
var week = 7 * day

// manages schedule times of the tasks
type timer struct {
	IsRepeating    bool           `json:"is_recurring"`
	Hour           int64          `json:"hour"`
	Minute         int64          `json:"minute"`
	Weekdays       []time.Weekday `json:"week_days"`
	MonthDay       int64          `json:"month_day"`
	OccurrenceType string         `json:"occurrence_type"`
	RunAt          int64          `json:"run_at"`
	StartAt        time.Time      `json:"start_at"`
	IsInfinite     bool           `json:"is_infinite"`
	UntilTime      time.Time      `json:"until_time"`
	RunCount       int64          `json:"run_count"`
	MaxRunCount    int64          `json:"max_run_count"`
}

func (t *timer) createNextDailyRunDate() {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), int(t.Hour), int(t.Minute), 0, 0, time.Local)

	for {
		if tm.Unix() > now.Unix() && tm.Unix() > t.StartAt.Unix() {
			break
		}

		tm = tm.AddDate(0, 0, 1)
	}

	t.OccurrenceType = "daily"
	t.RunAt = tm.Unix()
}

func (t *timer) createNextHourlyRunDate() {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), int(t.Minute), 0, 0, time.Local)

	for {
		if tm.Unix() > now.Unix() && tm.Unix() > t.StartAt.Unix() {
			break
		}

		tm = tm.Add(time.Hour)
	}

	t.OccurrenceType = "hourly"
	t.RunAt = tm.Unix()
}

func (t *timer) createNextWeeklyRunDay() {
	var weekdays = make(map[time.Weekday]bool)

	for _, weekday := range t.Weekdays {
		weekdays[weekday] = true
	}

	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), int(t.Hour), int(t.Minute), 0, 0, time.Local)

	for {
		if weekdays[tm.Weekday()] && tm.Unix() > t.StartAt.Unix() && tm.Unix() > now.Unix() {
			break
		}

		tm = tm.AddDate(0, 0, 1)
	}

	t.OccurrenceType = "week_days"
	t.RunAt = tm.Unix()
}

func (t *timer) createNextMonthlyRunDate() {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), int(t.Hour), int(t.Minute), 0, 0, time.Local)

	// Ignore months that have less days than specified
	for {
		start := time.Date(tm.Year(), tm.Month(), 1, 0, 0, 0, 0, time.Local)
		end := start.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

		if t.MonthDay <= int64(end.Day()) && tm.Unix() > t.StartAt.Unix() {
			tm = tm.AddDate(0, 0, int(t.MonthDay)-tm.Day())

			break
		}

		// Add month
		tm = tm.AddDate(0, 1, int(t.MonthDay)-tm.Day())
	}

	t.OccurrenceType = "monthly"
	t.RunAt = tm.Unix()
}

func (t *timer) setNextRun() {
	switch t.OccurrenceType {
	case "week_days":
		t.createNextWeeklyRunDay()
	case "daily":
		t.createNextDailyRunDate()
	case "hourly":
		t.createNextHourlyRunDate()
	case "monthly":
		t.createNextMonthlyRunDate()
	}
}

func (t *timer) shouldExecuteLater() (bool, time.Duration) {
	return t.StartAt.Unix() > time.Now().Unix(), t.StartAt.Sub(time.Now())
}

func (t *timer) isPastDate(tm time.Time) bool {
	return tm.Unix() < time.Now().Unix()
}
