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
	CurrentSession *Session  `json:"current_session"`
	Sessions       []Session `json:"sessions"`
	IsRunning      bool      `json:"is_running"`
	IsOnBreak      bool      `json:"is_on_break"`
	BreakStart     time.Time `json:"break_start"`
	SessionStart   time.Time `json:"session_start"`
	storage        *Storage  // Add storage field
}

func NewTimer() *Timer {
	return &Timer{
		Sessions: make([]Session, 0),
	}
}

// SetStorage sets the storage instance for the timer
func (t *Timer) SetStorage(s *Storage) {
	t.storage = s
}

func (t *Timer) Start() {
	if !t.IsRunning {
		now := time.Now()
		t.CurrentSession = &Session{
			Date: now.Format("2006-01-02"),
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

		t.CurrentSession.Duration = int64(time.Since(t.SessionStart).Seconds())
		// Only store sessions longer than 1 second
		if t.CurrentSession.Duration > 1 {
			t.Sessions = append(t.Sessions, *t.CurrentSession)
		}
		t.CurrentSession = nil
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
		breakDuration := time.Since(t.BreakStart)
		t.CurrentSession.BreakTime += int64(breakDuration.Seconds())
		t.IsOnBreak = false
	}
}

func (t *Timer) Reset() {
	// Don't save the current session when resetting
	t.CurrentSession = nil
	t.IsRunning = false
	t.IsOnBreak = false
}

func (t *Timer) GetCurrentTime() time.Duration {
	if !t.IsRunning || t.CurrentSession == nil {
		return 0
	}

	duration := time.Since(t.SessionStart)
	return duration - time.Duration(t.CurrentSession.BreakTime)*time.Second
}

func (t *Timer) GetCurrentBreakTime() time.Duration {
	if !t.IsRunning || t.CurrentSession == nil {
		return 0
	}

	if t.IsOnBreak {
		currentBreak := time.Since(t.BreakStart)
		return time.Duration(t.CurrentSession.BreakTime)*time.Second + currentBreak
	}
	return time.Duration(t.CurrentSession.BreakTime) * time.Second
}

func (t *Timer) GetWeeklyTime() time.Duration {
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

	// Add current session if running
	if t.IsRunning && t.CurrentSession != nil {
		sessionTime, err := time.Parse("2006-01-02", t.CurrentSession.Date)
		if err == nil { // Only add if date is valid
			sessionWeek := t.getWeekNumber(sessionTime)
			sessionYear := sessionTime.Year()

			if sessionWeek == currentWeek && sessionYear == currentYear {
				// Calculate current session duration excluding breaks
				currentDuration := time.Since(t.SessionStart)
				totalBreakTime := time.Duration(t.CurrentSession.BreakTime) * time.Second
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
		CurrentTime: t.GetCurrentTime(),
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
