package main

import (
	"log"

	"fyne.io/fyne/v2/app"
)

func main() {
	application := app.New()
	storage := NewStorage()

	// Load or create new timer
	timer, err := storage.LoadTimer()
	if err != nil {
		log.Printf("Error loading timer state: %v", err)
		timer = NewTimer()
	}

	// Set storage on timer for weekly calculations
	timer.SetStorage(storage)

	// Create UI
	ui := NewUI(application, timer, storage)

	// Set window to always be on top
	ui.window.SetFixedSize(true)

	// Show the window and start the application
	ui.Show()
	application.Run()
}
