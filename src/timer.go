package main

import (
	"encoding/json"
	"time"
)

type Session struct {
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	PauseTime time.Duration `json:"pause_time"`
}

type Timer struct {
	CurrentSession *Session  `json:"current_session"`
	Sessions       []Session `json:"sessions"`
	IsRunning      bool      `json:"is_running"`
	IsPaused       bool      `json:"is_paused"`
	PauseStart     time.Time `json:"pause_start"`
}

func NewTimer() *Timer {
	return &Timer{
		Sessions: make([]Session, 0),
	}
}

func (t *Timer) Start() {
	if !t.IsRunning {
		t.CurrentSession = &Session{
			StartTime: time.Now(),
		}
		t.IsRunning = true
	}
}

func (t *Timer) Stop() {
	if t.IsRunning {
		if t.IsPaused {
			t.StopPause()
		}

		t.CurrentSession.EndTime = time.Now()
		t.CurrentSession.Duration = t.CurrentSession.EndTime.Sub(t.CurrentSession.StartTime)
		t.Sessions = append(t.Sessions, *t.CurrentSession)
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
		t.CurrentSession.PauseTime += pauseDuration
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

	duration := time.Since(t.CurrentSession.StartTime)
	return duration - t.CurrentSession.PauseTime
}

func (t *Timer) GetCurrentPauseTime() time.Duration {
	if !t.IsRunning || t.CurrentSession == nil {
		return 0
	}

	if t.IsPaused {
		return t.CurrentSession.PauseTime + time.Since(t.PauseStart)
	}
	return t.CurrentSession.PauseTime
}

func (t *Timer) GetWeeklyTime() time.Duration {
	var total time.Duration
	currentWeek := t.getWeekNumber(time.Now())
	currentYear := time.Now().Year()

	// Sum up completed sessions from this week
	for _, session := range t.Sessions {
		sessionWeek := t.getWeekNumber(session.StartTime)
		sessionYear := session.StartTime.Year()

		if sessionWeek == currentWeek && sessionYear == currentYear {
			total += session.Duration - session.PauseTime
		}
	}

	// Add current session if running
	if t.IsRunning && t.CurrentSession != nil {
		sessionWeek := t.getWeekNumber(t.CurrentSession.StartTime)
		sessionYear := t.CurrentSession.StartTime.Year()

		if sessionWeek == currentWeek && sessionYear == currentYear {
			total += t.GetCurrentTime()
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
