package main

import (
	"os"
	"testing"
	"time"

	timetracker "github.com/LitschiW/timetracker/pkg/timetracker"
)

func TestTimerBasicOperations(t *testing.T) {
	timer := timetracker.NewTimer()

	// Test initial state
	if timer.IsRunning {
		t.Error("New timer should not be running")
	}
	if timer.IsOnBreak {
		t.Error("New timer should not be on break")
	}
	if timer.GetCurrentTime() != 0 {
		t.Error("New timer should have 0 current time")
	}

	// Test start
	timer.Start()
	time.Sleep(2 * time.Second)
	if !timer.IsRunning {
		t.Error("Timer should be running after Start")
	}
	if timer.CurrentSession == nil {
		t.Error("Current session should be initialized after Start")
	}

	// Test break
	timer.StartBreak()
	time.Sleep(1 * time.Second)
	timer.StopBreak()
	if timer.IsOnBreak {
		t.Error("Timer should not be on break after StopBreak")
	}
	if timer.CurrentSession.BreakTime == 0 {
		t.Error("Break time should be recorded")
	}

	// Test stop
	timer.Stop()
	if timer.IsRunning {
		t.Error("Timer should not be running after Stop")
	}
	if len(timer.Sessions) != 1 {
		t.Error("Session should be added to Sessions after Stop")
	}
}

func TestTimerWeeklyTotal(t *testing.T) {
	// Create temporary files for testing
	tmpJSON := "test_current_session.json"
	tmpCSV := "test_sessions.csv"
	defer os.Remove(tmpJSON)
	defer os.Remove(tmpCSV)

	storage := timetracker.NewStorage()

	timer := timetracker.NewTimer()
	timer.SetStorage(storage)

	// Test empty state
	if total := timer.GetWeeklyTime(); total != 0 {
		t.Errorf("Empty timer should have 0 weekly total, got %v", total)
	}

	// Add a session
	timer.Start()
	time.Sleep(2 * time.Second)
	timer.Stop()

	// Save to storage
	if err := storage.SaveTimer(timer); err != nil {
		t.Fatalf("Failed to save timer: %v", err)
	}

	// Load new timer from storage
	newTimer, err := storage.LoadTimer()
	if err != nil {
		t.Fatalf("Failed to load timer: %v", err)
	}

	// Check weekly total
	total := newTimer.GetWeeklyTime()
	if total < 2*time.Second {
		t.Errorf("Weekly total should be at least 2 seconds, got %v", total)
	}
}

func TestBreakTimeExclusion(t *testing.T) {
	timer := timetracker.NewTimer()

	// Start a session
	timer.Start()
	time.Sleep(3 * time.Second)

	// Take a break
	timer.StartBreak()
	time.Sleep(1 * time.Second)
	timer.StopBreak()

	// Continue working
	time.Sleep(2 * time.Second)
	timer.Stop()

	// Total time should be ~5 seconds (3 + 2), not 6 seconds
	session := timer.Sessions[0]
	workTime := session.Duration - session.BreakTime
	if workTime < 4 || workTime > 5 {
		t.Errorf("Work time should be ~5 seconds, got %v seconds", workTime)
	}
}

func TestShortSessionExclusion(t *testing.T) {
	timer := timetracker.NewTimer()

	// Start and immediately stop
	timer.Start()
	timer.Stop()

	// Session less than 1 second should not be recorded
	if len(timer.Sessions) != 0 {
		t.Error("Short session should not be recorded")
	}
}

func TestWeeklyTotalCache(t *testing.T) {
	// Create temporary files for testing
	tmpJSON := "test_current_session.json"
	tmpCSV := "test_sessions.csv"
	defer os.Remove(tmpJSON)
	defer os.Remove(tmpCSV)

	storage := timetracker.NewStorage()

	timer := timetracker.NewTimer()
	timer.SetStorage(storage)

	// Add multiple sessions
	for i := 0; i < 3; i++ {
		timer.Start()
		time.Sleep(2 * time.Second)
		timer.Stop()
		if err := storage.SaveTimer(timer); err != nil {
			t.Fatalf("Failed to save timer: %v", err)
		}
	}

	// Load timer and check cached total
	newTimer, err := storage.LoadTimer()
	if err != nil {
		t.Fatalf("Failed to load timer: %v", err)
	}

	// Weekly total should be cached and not require CSV reads
	total1 := newTimer.GetWeeklyTime()
	total2 := newTimer.GetWeeklyTime()

	if total1 != total2 {
		t.Error("Cached weekly totals should be equal")
	}
	if total1 < 6*time.Second {
		t.Errorf("Weekly total should be at least 6 seconds, got %v", total1)
	}
}
