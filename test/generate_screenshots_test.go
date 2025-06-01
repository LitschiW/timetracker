package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	timetracker "github.com/LitschiW/timetracker/pkg/timetracker"
)

func TestGenerateScreenshots(t *testing.T) {
	// Create or clear screenshots directory
	screenshotDir := "screenshots"
	if err := os.RemoveAll(screenshotDir); err != nil {
		t.Fatalf("Failed to clear screenshots directory: %v", err)
	}
	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		t.Fatalf("Failed to create screenshots directory: %v", err)
	}

	// Create a new test app
	a := test.NewApp()
	storage := timetracker.NewStorage()
	timer := timetracker.NewTimer()
	ui := timetracker.NewUI(a, timer, storage)

	// Show the window
	ui.Window.Resize(fyne.NewSize(timetracker.WindowWidth, timetracker.WindowHeight))
	w := test.NewWindow(ui.Window.Content())
	w.Resize(fyne.NewSize(timetracker.WindowWidth, timetracker.WindowHeight))

	// Function to capture and save screenshot
	captureScreen := func(name string) {
		time.Sleep(500 * time.Millisecond) // Wait for UI to update

		// Capture the window content
		c := w.Canvas()
		img := c.Capture()

		// Save the image
		path := filepath.Join(screenshotDir, fmt.Sprintf("%s.png", name))
		f, err := os.Create(path)
		if err != nil {
			t.Errorf("Failed to create file: %v", err)
			return
		}
		defer f.Close()

		if err := png.Encode(f, img); err != nil {
			t.Errorf("Failed to encode image: %v", err)
			return
		}

		t.Logf("Saved screenshot: %s", path)
	}

	// 1. Initial state
	captureScreen("initial_state")

	// 2. Working session started
	ui.HandleStartStop()
	captureScreen("working_session")

	// 3. Break started (wait for 3 dots)
	ui.HandleBreak()
	time.Sleep(4500 * time.Millisecond) // Wait for 3 dots to appear
	captureScreen("break_state")

	// Cleanup
	w.Close()
}
