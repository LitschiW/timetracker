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
	breakLabel   *widget.Label
	weeklyLabel  *widget.Label
	startButton  *widget.Button
	breakButton  *widget.Button
	cancelButton *widget.Button
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
	ui.timeLabel = widget.NewLabel("Current Session: 0:00:00")
	ui.breakLabel = widget.NewLabel("Current Break Time: 0:00:00")
	ui.weeklyLabel = widget.NewLabel("This Weeks Total: 0:00:00")

	ui.startButton = widget.NewButton("Start Working Session", ui.handleStartStop)
	ui.breakButton = widget.NewButton("Start Break", ui.handleBreak)
	ui.cancelButton = widget.NewButton("Cancel Working Session", ui.handleCancel)

	// Update button states based on timer state
	if ui.timer.IsRunning {
		ui.startButton.SetText("Stop Working Session")
	}
	if ui.timer.IsOnBreak {
		ui.breakButton.SetText("Resume Working Session")
	}
}

func (ui *UI) layoutWidgets() {
	buttons := container.NewVBox(
		layout.NewSpacer(),
		ui.startButton,
		ui.breakButton,
		ui.cancelButton,
		layout.NewSpacer(),
	)

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), ui.timeLabel, layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(), ui.breakLabel, layout.NewSpacer()),
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
	fyne.Do(
		func() {
			ui.timeLabel.SetText("Current Session: " + ui.formatDuration(ui.timer.GetCurrentTime()))
			ui.breakLabel.SetText("Current Break Time: " + ui.formatDuration(ui.timer.GetCurrentBreakTime()))
			ui.weeklyLabel.SetText("This Weeks Total: " + ui.formatDuration(ui.timer.GetWeeklyTime()))
		})
}

func (ui *UI) handleStartStop() {
	if ui.timer.IsRunning {
		ui.timer.Stop()
		ui.startButton.SetText("Start Working Session")
		ui.breakButton.SetText("Start Break")
	} else {
		ui.timer.Start()
		ui.startButton.SetText("Stop Working Session")
	}
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) handleBreak() {
	if ui.timer.IsOnBreak || !ui.timer.IsRunning {
		ui.timer.StopBreak()
		ui.breakButton.SetText("Start Break")
	} else {
		ui.timer.StartBreak()
		ui.breakButton.SetText("Stop Break")
	}
	ui.storage.SaveTimer(ui.timer)
}

func (ui *UI) handleCancel() {
	ui.timer.Reset()
	ui.startButton.SetText("Start Working Session")
	ui.breakButton.SetText("Start Break")
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
