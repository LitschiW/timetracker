package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Storage struct {
	filename string
}

func NewStorage() *Storage {
	return &Storage{
		filename: "timetracker.json",
	}
}

func (s *Storage) SaveTimer(timer *Timer) error {
	data, err := json.Marshal(timer)
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *Storage) LoadTimer() (*Timer, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return NewTimer(), nil
		}
		return nil, err
	}

	var timer Timer
	if err := json.Unmarshal(data, &timer); err != nil {
		return nil, err
	}

	return &timer, nil
}
