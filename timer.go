package main

import (
	"encoding/json"
	"time"
)

type Session struct {
	Date      string `json:"date"`
	Duration  int64  `json:"duration_s"`
	PauseTime int64  `json:"pause_time_s"`
}

type Timer struct {
	CurrentSession *Session  `json:"current_session"`
	Sessions       []Session `json:"sessions"`
	IsRunning      bool      `json:"is_running"`
	IsPaused       bool      `json:"is_paused"`
	PauseStart     time.Time `json:"pause_start"`
	SessionStart   time.Time `json:"session_start"`
}

func NewTimer() *Timer {
	return &Timer{
		Sessions: make([]Session, 0),
	}
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
		if t.IsPaused {
			t.StopPause()
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

func (t *Timer) StartPause() {
	if t.IsRunning && !t.IsPaused {
		t.PauseStart = time.Now()
		t.IsPaused = true
	}
}

func (t *Timer) StopPause() {
	if t.IsPaused {
		pauseDuration := time.Since(t.PauseStart)
		t.CurrentSession.PauseTime += int64(pauseDuration.Seconds())
		t.IsPaused = false
	}
}

func (t *Timer) Reset() {
	if t.IsRunning {
		t.Stop()
	}
	t.CurrentSession = nil
	t.IsRunning = false
	t.IsPaused = false
}

func (t *Timer) GetCurrentTime() time.Duration {
	if !t.IsRunning || t.CurrentSession == nil {
		return 0
	}

	duration := time.Since(t.SessionStart)
	return duration - time.Duration(t.CurrentSession.PauseTime)*time.Second
}

func (t *Timer) GetCurrentPauseTime() time.Duration {
	if !t.IsRunning || t.CurrentSession == nil {
		return 0
	}

	if t.IsPaused {
		currentPause := time.Since(t.PauseStart)
		return time.Duration(t.CurrentSession.PauseTime)*time.Second + currentPause
	}
	return time.Duration(t.CurrentSession.PauseTime) * time.Second
}

func (t *Timer) GetWeeklyTime() time.Duration {
	var total time.Duration
	currentWeek := t.getWeekNumber(time.Now())
	currentYear := time.Now().Year()

	// Sum up completed sessions from this week
	for _, session := range t.Sessions {
		sessionTime, err := time.Parse("2006-01-02", session.Date)
		if err != nil {
			continue // Skip invalid dates
		}
		sessionWeek := t.getWeekNumber(sessionTime)
		sessionYear := sessionTime.Year()

		if sessionWeek == currentWeek && sessionYear == currentYear {
			durationDiff := session.Duration - session.PauseTime
			total += time.Duration(durationDiff) * time.Second
		}
	}

	// Add current session if running
	if t.IsRunning && t.CurrentSession != nil {
		sessionTime, err := time.Parse("2006-01-02", t.CurrentSession.Date)
		if err == nil { // Only add if date is valid
			sessionWeek := t.getWeekNumber(sessionTime)
			sessionYear := sessionTime.Year()

			if sessionWeek == currentWeek && sessionYear == currentYear {
				total += t.GetCurrentTime()
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
