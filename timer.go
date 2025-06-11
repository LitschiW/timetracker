package main

import (
	"encoding/json"
	"time"
)

type Session struct {
	Date      string `json:"date"`
	Duration  int64  `json:"duration_s"`
	BreakTime int64  `json:"break_time_s"`
}

type Timer struct {
	TodaySession        *Session      `json:"today_session"`
	Sessions            []Session     `json:"sessions"`
	IsRunning           bool          `json:"is_running"`
	IsOnBreak           bool          `json:"is_on_break"`
	BreakStart          time.Time     `json:"break_start"`
	SessionStart        time.Time     `json:"session_start"`
	DayFirstStart       time.Time     `json:"day_first_start"`
	DailyTotal          time.Duration `json:"daily_total"`
	YesterdayTotal      time.Duration `json:"yesterday_total"`
	YesterdayFirstStart time.Time     `json:"yesterday_first_start"`
	storage             *Storage
	weeklyTotal         time.Duration
}

func NewTimer() *Timer {
	return &Timer{
		Sessions: make([]Session, 0),
	}
}

// SetStorage sets the storage instance for the timer and updates weekly total
func (t *Timer) SetStorage(s *Storage) {
	t.storage = s
	t.updateWeeklyTotal() // Initialize weekly total on storage set
}

func (t *Timer) Start() {
	if !t.IsRunning {
		now := time.Now()
		currentDate := now.Format("2006-01-02")

		// Check for day transition
		t.checkAndHandleDayTransition()

		// Check if this is the first session of a new day
		if t.DayFirstStart.IsZero() {
			t.DayFirstStart = now
			t.DailyTotal = 0
		}

		t.TodaySession = &Session{
			Date: currentDate,
		}
		t.SessionStart = now
		t.IsRunning = true
	}
}

func (t *Timer) Stop() {
	if t.IsRunning {
		if t.IsOnBreak {
			t.StopBreak()
		}

		sessionDuration := time.Since(t.SessionStart)
		breakDuration := time.Duration(t.TodaySession.BreakTime) * time.Second
		workDuration := sessionDuration - breakDuration

		t.TodaySession.Duration = int64(sessionDuration.Seconds())
		// Only store sessions longer than 1 second
		if t.TodaySession.Duration > 1 {
			t.Sessions = append(t.Sessions, *t.TodaySession)
			t.updateWeeklyTotal() // Update weekly total when adding new session

			// Update daily total
			if workDuration > 0 {
				t.DailyTotal += workDuration
			}
		}
		t.TodaySession = nil
		t.IsRunning = false
	}
}

func (t *Timer) StartBreak() {
	if t.IsRunning && !t.IsOnBreak {
		t.BreakStart = time.Now()
		t.IsOnBreak = true
	}
}

func (t *Timer) StopBreak() {
	if t.IsOnBreak {
		breakDuration := time.Since(t.BreakStart).Seconds()
		t.TodaySession.BreakTime += int64(breakDuration)
		t.IsOnBreak = false
	}
}

func (t *Timer) Reset() {
	// Don't save the current session when resetting
	t.TodaySession = nil
	t.IsRunning = false
	t.IsOnBreak = false
}

func (t *Timer) GetTodaySessionTime() time.Duration {
	if !t.IsRunning || t.TodaySession == nil {
		return 0
	}

	duration := time.Since(t.SessionStart)
	return duration - time.Duration(t.TodaySession.BreakTime)*time.Second
}

func (t *Timer) GetCurrentBreakTime() time.Duration {
	if !t.IsRunning || t.TodaySession == nil {
		return 0
	}

	if t.IsOnBreak {
		currentBreak := time.Since(t.BreakStart)
		return time.Duration(t.TodaySession.BreakTime)*time.Second + currentBreak
	}
	return time.Duration(t.TodaySession.BreakTime) * time.Second
}

// updateWeeklyTotal recalculates and caches the weekly total
func (t *Timer) updateWeeklyTotal() {
	var total time.Duration
	currentWeek := t.getWeekNumber(time.Now())
	currentYear := time.Now().Year()

	// Load historical sessions from CSV
	if t.storage != nil {
		if sessions, err := t.storage.loadSessionsFromCSV(); err == nil {
			// Sum up completed sessions from this week
			for _, session := range sessions {
				sessionTime, err := time.Parse("2006-01-02", session.Date)
				if err != nil {
					continue // Skip invalid dates
				}
				sessionWeek := t.getWeekNumber(sessionTime)
				sessionYear := sessionTime.Year()

				if sessionWeek == currentWeek && sessionYear == currentYear {
					durationDiff := session.Duration - session.BreakTime
					if durationDiff > 0 { // Only count positive durations
						total += time.Duration(durationDiff) * time.Second
					}
				}
			}
		}
	}

	// Add completed sessions from memory that haven't been saved yet
	for _, session := range t.Sessions {
		sessionTime, err := time.Parse("2006-01-02", session.Date)
		if err != nil {
			continue // Skip invalid dates
		}
		sessionWeek := t.getWeekNumber(sessionTime)
		sessionYear := sessionTime.Year()

		if sessionWeek == currentWeek && sessionYear == currentYear {
			durationDiff := session.Duration - session.BreakTime
			if durationDiff > 0 { // Only count positive durations
				total += time.Duration(durationDiff) * time.Second
			}
		}
	}

	t.weeklyTotal = total
}

func (t *Timer) GetWeeklyTime() time.Duration {
	// Return cached total plus current session if running
	total := t.weeklyTotal

	// Add current session if running
	if t.IsRunning && t.TodaySession != nil {
		sessionTime, err := time.Parse("2006-01-02", t.TodaySession.Date)
		if err == nil { // Only add if date is valid
			sessionWeek := t.getWeekNumber(sessionTime)
			sessionYear := sessionTime.Year()
			currentWeek := t.getWeekNumber(time.Now())
			currentYear := time.Now().Year()

			if sessionWeek == currentWeek && sessionYear == currentYear {
				// Calculate current session duration excluding breaks
				currentDuration := time.Since(t.SessionStart)
				totalBreakTime := time.Duration(t.TodaySession.BreakTime) * time.Second
				if t.IsOnBreak {
					totalBreakTime += time.Since(t.BreakStart)
				}
				workDuration := currentDuration - totalBreakTime
				if workDuration > 0 {
					total += workDuration
				}
			}
		}
	}

	return total
}

func (t *Timer) getWeekNumber(date time.Time) int {
	_, week := date.ISOWeek()
	return week
}

func (t *Timer) MarshalJSON() ([]byte, error) {
	type Alias Timer
	return json.Marshal(&struct {
		*Alias
		CurrentTime time.Duration `json:"current_time"`
		WeeklyTime  time.Duration `json:"weekly_time"`
	}{
		Alias:       (*Alias)(t),
		CurrentTime: t.GetTodaySessionTime(),
		WeeklyTime:  t.GetWeeklyTime(),
	})
}

func (t *Timer) UnmarshalJSON(data []byte) error {
	type Alias Timer
	aux := &struct {
		*Alias
		CurrentTime time.Duration `json:"current_time"`
		WeeklyTime  time.Duration `json:"weekly_time"`
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// Add method to get current day's total time
func (t *Timer) GetDailyTime() time.Duration {
	total := t.DailyTotal

	// Add current session if running
	if t.IsRunning && t.TodaySession != nil {
		currentDuration := time.Since(t.SessionStart)
		totalBreakTime := time.Duration(t.TodaySession.BreakTime) * time.Second
		if t.IsOnBreak {
			totalBreakTime += time.Since(t.BreakStart)
		}
		workDuration := currentDuration - totalBreakTime
		if workDuration > 0 {
			total += workDuration
		}
	}

	return total
}

// Add method to get formatted first start time of the day
func (t *Timer) GetDayFirstStartTime() string {
	if t.DayFirstStart.IsZero() {
		return "Not started today"
	}
	return t.DayFirstStart.Format("15:04:05")
}

// Add method to handle day transition
func (t *Timer) checkAndHandleDayTransition() {
	now := time.Now()
	if !t.DayFirstStart.IsZero() && t.DayFirstStart.Format("2006-01-02") != now.Format("2006-01-02") {
		// Store yesterday's data before resetting
		t.YesterdayTotal = t.DailyTotal
		t.YesterdayFirstStart = t.DayFirstStart

		// Reset today's tracking
		t.DayFirstStart = time.Time{}
		t.DailyTotal = 0
	}
}

// Add method to get yesterday's first start time
func (t *Timer) GetYesterdayFirstStartTime() string {
	if t.YesterdayFirstStart.IsZero() {
		return "No data"
	}
	return t.YesterdayFirstStart.Format("15:04:05")
}
