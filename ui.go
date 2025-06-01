package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	windowTitle  = "Time Tracker"
	windowWidth  = 340
	windowHeight = 250

	// Button text
	textStart      = "▶️ Start Working Session"
	textStop       = "💾 Stop & Save Session"
	textStartBreak = "☕ Start Break"
	textStopBreak  = "☕ Stop Break"
	textCancel     = "❌ Cancel Working Session"

	// Label formats
	formatCurrentSession = "Current Session: %s"
	formatBreakTime      = "Current Break Time: %s"
	formatWeeklyTotal    = "This Week's Total: %s"
)

type UI struct {
	window       fyne.Window
	timer        *Timer
	storage      *Storage
	timeLabel    *widget.Label
	breakLabel   *widget.Label
	weeklyLabel  *widget.Label
	timeDesc     *widget.Label
	breakDesc    *widget.Label
	weeklyDesc   *widget.Label
	startButton  *widget.Button
	breakButton  *widget.Button
	cancelButton *widget.Button
	updateTicker *time.Ticker
	quitChan     chan struct{}
}

func NewUI(app fyne.App, timer *Timer, storage *Storage) *UI {
	ui := &UI{
		window:   app.NewWindow(windowTitle),
		timer:    timer,
		storage:  storage,
		quitChan: make(chan struct{}),
	}

	ui.window.SetOnClosed(ui.handleClose)
	ui.window.Resize(fyne.NewSize(windowWidth, windowHeight))
	ui.window.SetFixedSize(true)
	ui.window.CenterOnScreen()

	ui.createWidgets()
	ui.layoutWidgets()
	ui.startUpdateTicker()

	return ui
}

func (ui *UI) createWidgets() {
	// Create labels with centered text
	ui.timeLabel = widget.NewLabel("0:00:00")
	ui.breakLabel = widget.NewLabel("0:00:00")
	ui.weeklyLabel = widget.NewLabel("0:00:00")

	// Create static label descriptions
	ui.timeDesc = widget.NewLabel("Current Session:")
	ui.breakDesc = widget.NewLabel("Current Break Time:")
	ui.weeklyDesc = widget.NewLabel("This Week's Total:")

	// Set alignment for description labels (right-aligned)
	ui.timeDesc.Alignment = fyne.TextAlignTrailing
	ui.breakDesc.Alignment = fyne.TextAlignTrailing
	ui.weeklyDesc.Alignment = fyne.TextAlignTrailing

	// Set alignment for value labels (centered)
	ui.timeLabel.Alignment = fyne.TextAlignCenter
	ui.breakLabel.Alignment = fyne.TextAlignCenter
	ui.weeklyLabel.Alignment = fyne.TextAlignCenter

	// Create buttons without theme icons since we're using Unicode icons
	ui.startButton = widget.NewButton(textStart, ui.handleStartStop)
	ui.breakButton = widget.NewButton(textStartBreak, ui.handleBreak)
	ui.cancelButton = widget.NewButton(textCancel, ui.handleCancel)

	// Initialize button states
	ui.updateButtonStates()
}

func (ui *UI) layoutWidgets() {
	// Create grid for labels
	labelGrid := container.NewGridWithColumns(2,
		ui.timeDesc, ui.timeLabel,
		ui.breakDesc, ui.breakLabel,
		ui.weeklyDesc, ui.weeklyLabel,
	)

	// Create button container
	buttons := container.NewVBox(
		layout.NewSpacer(),
		ui.startButton,
		layout.NewSpacer(),
		ui.breakButton,
		layout.NewSpacer(),
		ui.cancelButton,
		layout.NewSpacer(),
	)

	// Main content layout with centered elements
	content := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), labelGrid, layout.NewSpacer()),
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), buttons, layout.NewSpacer()),
	)

	ui.window.SetContent(content)
}

func (ui *UI) formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

func (ui *UI) updateLabels() {
	fyne.Do(
		func() {
			ui.timeLabel.SetText(ui.formatDuration(ui.timer.GetCurrentTime()))
			ui.breakLabel.SetText(ui.formatDuration(ui.timer.GetCurrentBreakTime()))
			ui.weeklyLabel.SetText(ui.formatDuration(ui.timer.GetWeeklyTime()))

			// Update break button text with animated dots when on break
			if ui.timer.IsOnBreak {
				dots := strings.Repeat(".", int(time.Now().Unix()%4))
				ui.breakButton.SetText("☕ Stop Break" + dots)
			} else if ui.timer.IsRunning {
				ui.breakButton.SetText(textStartBreak)
			}
		})
}

func (ui *UI) updateButtonStates() {
	// Start/Stop button text
	if ui.timer.IsRunning {
		ui.startButton.SetText(textStop)
		ui.cancelButton.Enable()
		ui.breakButton.Enable()
	} else {
		ui.startButton.SetText(textStart)
		ui.cancelButton.Disable()
		ui.breakButton.Disable()
	}

	// Break button text
	if ui.timer.IsOnBreak {
		ui.breakButton.SetText(textStopBreak)
	} else {
		ui.breakButton.SetText(textStartBreak)
	}
}

func (ui *UI) handleStartStop() {
	if ui.timer.IsRunning {
		ui.timer.Stop()
	} else {
		ui.timer.Start()
	}
	ui.updateButtonStates()
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) handleBreak() {
	if ui.timer.IsOnBreak || !ui.timer.IsRunning {
		ui.timer.StopBreak()
	} else {
		ui.timer.StartBreak()
	}
	ui.updateButtonStates()
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) handleCancel() {
	ui.timer.Reset()
	ui.updateButtonStates()
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) startUpdateTicker() {
	ui.updateTicker = time.NewTicker(250 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ui.updateTicker.C:
				ui.updateLabels()
			case <-ui.quitChan:
				return
			}
		}
	}()
}

func (ui *UI) handleClose() {
	ui.updateTicker.Stop()
	close(ui.quitChan)
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) Show() {
	ui.window.Show()
}
