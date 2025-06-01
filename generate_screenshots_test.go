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
	storage := NewStorage()
	timer := NewTimer()
	ui := NewUI(a, timer, storage)

	// Show the window
	ui.window.Resize(fyne.NewSize(windowWidth, windowHeight))
	w := test.NewWindow(ui.window.Content())
	w.Resize(fyne.NewSize(windowWidth, windowHeight))

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

	// Test different states

	// 1. Initial state
	captureScreen("initial_state")

	// 2. Working session started
	ui.handleStartStop()
	captureScreen("working_session")

	// 3. Break started (wait for 3 dots)
	ui.handleBreak()
	time.Sleep(3500 * time.Millisecond) // Wait for 3 dots to appear
	captureScreen("break_state")

	// 4. Break stopped
	ui.handleBreak()
	captureScreen("break_stopped")

	// 5. Session stopped
	ui.handleStartStop()
	captureScreen("session_stopped")

	// 6. Session cancelled
	ui.handleStartStop()
	ui.handleCancel()
	captureScreen("session_cancelled")

	// Cleanup
	w.Close()
}
