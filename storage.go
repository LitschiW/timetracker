package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Storage struct {
	jsonFile string
	csvFile  string
}

func NewStorage() *Storage {
	return &Storage{
		jsonFile: "current_session.json",
		csvFile:  "sessions.csv",
	}
}

func (s *Storage) SaveTimer(timer *Timer) error {
	// First save completed sessions to CSV if any exist
	if len(timer.Sessions) > 0 {
		if err := s.appendSessionsToCSV(timer.Sessions); err != nil {
			return fmt.Errorf("failed to save sessions to CSV: %w", err)
		}
		// Clear sessions from timer after saving to CSV
		timer.Sessions = []Session{}
		// Update weekly total after saving to CSV
		timer.updateWeeklyTotal()
	}

	// Then save current state to JSON
	data, err := json.MarshalIndent(timer, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal timer: %w", err)
	}

	return os.WriteFile(s.jsonFile, data, 0644)
}

func (s *Storage) LoadTimer() (*Timer, error) {
	// Load current state from JSON
	data, err := os.ReadFile(s.jsonFile)
	if err != nil {
		if os.IsNotExist(err) {
			timer := NewTimer()
			timer.SetStorage(s)
			return timer, nil
		}
		return nil, err
	}

	var timer Timer
	if err := json.Unmarshal(data, &timer); err != nil {
		return nil, err
	}

	// Initialize empty sessions slice
	timer.Sessions = make([]Session, 0)

	// Check if we need to handle day transition
	now := time.Now()
	if !timer.DayFirstStart.IsZero() && timer.DayFirstStart.Format("2006-01-02") != now.Format("2006-01-02") {
		// Store yesterday's data
		timer.YesterdayTotal = timer.DailyTotal
		timer.YesterdayFirstStart = timer.DayFirstStart

		// Reset today's tracking
		timer.DayFirstStart = time.Time{}
		timer.DailyTotal = 0
	}

	// Set storage and update weekly total
	timer.SetStorage(s)

	return &timer, nil
}

func (s *Storage) appendSessionsToCSV(sessions []Session) error {
	// Open file in append mode, create if not exists
	file, err := os.OpenFile(s.csvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// If file is empty, write header
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() == 0 {
		if err := writer.Write([]string{"date", "duration_s", "break_time_s"}); err != nil {
			return err
		}
	}

	// Write sessions
	for _, session := range sessions {
		record := []string{
			session.Date,
			strconv.FormatInt(session.Duration, 10),
			strconv.FormatInt(session.BreakTime, 10),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) loadSessionsFromCSV() ([]Session, error) {
	// If file doesn't exist, return empty slice
	if _, err := os.Stat(s.csvFile); os.IsNotExist(err) {
		return []Session{}, nil
	}

	file, err := os.Open(s.csvFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and skip header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var sessions []Session
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if len(record) != 3 {
			continue // Skip invalid records
		}

		duration, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		breakTime, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			continue
		}

		sessions = append(sessions, Session{
			Date:      record[0],
			Duration:  duration,
			BreakTime: breakTime,
		})
	}

	return sessions, nil
}
