package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	window       fyne.Window
	timer        *Timer
	storage      *Storage
	timeLabel    *widget.Label
	pauseLabel   *widget.Label
	weeklyLabel  *widget.Label
	startButton  *widget.Button
	pauseButton  *widget.Button
	resetButton  *widget.Button
	updateTicker *time.Ticker
	quitChan     chan struct{}
}

func NewUI(app fyne.App, timer *Timer, storage *Storage) *UI {
	ui := &UI{
		window:   app.NewWindow("Time Tracker"),
		timer:    timer,
		storage:  storage,
		quitChan: make(chan struct{}),
	}

	ui.createWidgets()
	ui.layoutWidgets()
	ui.startUpdateTicker()

	ui.window.SetOnClosed(ui.handleClose)
	ui.window.Resize(fyne.NewSize(300, 150))
	ui.window.SetFixedSize(true)
	ui.window.CenterOnScreen()

	return ui
}

func (ui *UI) createWidgets() {
	ui.timeLabel = widget.NewLabel("Time: 0:00:00")
	ui.pauseLabel = widget.NewLabel("Pause: 0:00:00")
	ui.weeklyLabel = widget.NewLabel("Weekly: 0:00:00")

	ui.startButton = widget.NewButton("Start", ui.handleStartStop)
	ui.pauseButton = widget.NewButton("Pause", ui.handlePause)
	ui.resetButton = widget.NewButton("Reset", ui.handleReset)

	// Update button states based on timer state
	if ui.timer.IsRunning {
		ui.startButton.SetText("Stop")
	}
	if ui.timer.IsPaused {
		ui.pauseButton.SetText("Resume")
	}
}

func (ui *UI) layoutWidgets() {
	buttons := container.NewHBox(
		layout.NewSpacer(),
		ui.startButton,
		ui.pauseButton,
		ui.resetButton,
		layout.NewSpacer(),
	)

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), ui.timeLabel, layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(), ui.pauseLabel, layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(), ui.weeklyLabel, layout.NewSpacer()),
		buttons,
		layout.NewSpacer(),
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
	ui.timeLabel.SetText("Time: " + ui.formatDuration(ui.timer.GetCurrentTime()))
	ui.pauseLabel.SetText("Pause: " + ui.formatDuration(ui.timer.GetCurrentPauseTime()))
	ui.weeklyLabel.SetText("Weekly: " + ui.formatDuration(ui.timer.GetWeeklyTime()))
}

func (ui *UI) handleStartStop() {
	if ui.timer.IsRunning {
		ui.timer.Stop()
		ui.startButton.SetText("Start")
	} else {
		ui.timer.Start()
		ui.startButton.SetText("Stop")
	}
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) handlePause() {
	if ui.timer.IsPaused {
		ui.timer.StopPause()
		ui.pauseButton.SetText("Pause")
	} else {
		ui.timer.StartPause()
		ui.pauseButton.SetText("Resume")
	}
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) handleReset() {
	ui.timer.Reset()
	ui.startButton.SetText("Start")
	ui.pauseButton.SetText("Pause")
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) startUpdateTicker() {
	ui.updateTicker = time.NewTicker(time.Second)
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
