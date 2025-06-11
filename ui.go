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
	windowWidth  = 350
	windowHeight = 300

	// Button text
	textStart      = "▶️ Start Working Session"
	textStop       = "💾 Stop & Save Session"
	textStartBreak = "☕ Start Break"
	textStopBreak  = "☕ Stop Break"
	textCancel     = "❌ Cancel Working Session"

	// Label formats
	formatTodaySession   = "Today's Session: "
	formatBreakTime      = "Current Break Time: "
	formatWeeklyTotal    = "This Week's Total: "
	formatDailyTotal     = "Today's Total: "
	formatFirstStart     = "Started at: "
	formatYesterdayStats = "Yesterday's Stats"
)

type UI struct {
	window              fyne.Window
	timer               *Timer
	storage             *Storage
	todayTimeLabel      *widget.Label
	breakLabel          *widget.Label
	weeklyLabel         *widget.Label
	dailyLabel          *widget.Label
	firstStartLabel     *widget.Label
	yesterdayDailyLabel *widget.Label
	yesterdayStartLabel *widget.Label
	todayTimeDesc       *widget.Label
	breakDesc           *widget.Label
	weeklyDesc          *widget.Label
	dailyDesc           *widget.Label
	firstStartDesc      *widget.Label
	yesterdayTitle      *widget.Label
	yesterdayDailyDesc  *widget.Label
	yesterdayStartDesc  *widget.Label
	startButton         *widget.Button
	breakButton         *widget.Button
	cancelButton        *widget.Button
	updateTicker        *time.Ticker
	quitChan            chan struct{}
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
	// Create labels with proper alignment
	ui.todayTimeDesc = widget.NewLabelWithStyle(formatTodaySession, fyne.TextAlignTrailing, fyne.TextStyle{})
	ui.todayTimeLabel = widget.NewLabelWithStyle("0:00:00", fyne.TextAlignLeading, fyne.TextStyle{})
	ui.breakDesc = widget.NewLabelWithStyle(formatBreakTime, fyne.TextAlignTrailing, fyne.TextStyle{})
	ui.breakLabel = widget.NewLabelWithStyle("0:00:00", fyne.TextAlignLeading, fyne.TextStyle{})
	ui.weeklyDesc = widget.NewLabelWithStyle(formatWeeklyTotal, fyne.TextAlignTrailing, fyne.TextStyle{})
	ui.weeklyLabel = widget.NewLabelWithStyle("0:00:00", fyne.TextAlignLeading, fyne.TextStyle{})
	ui.dailyDesc = widget.NewLabelWithStyle(formatDailyTotal, fyne.TextAlignTrailing, fyne.TextStyle{})
	ui.dailyLabel = widget.NewLabelWithStyle("0:00:00", fyne.TextAlignLeading, fyne.TextStyle{})
	ui.firstStartDesc = widget.NewLabelWithStyle(formatFirstStart, fyne.TextAlignTrailing, fyne.TextStyle{})
	ui.firstStartLabel = widget.NewLabelWithStyle("Not started today", fyne.TextAlignLeading, fyne.TextStyle{})

	// Create yesterday's stats labels with proper alignment
	ui.yesterdayTitle = widget.NewLabelWithStyle(formatYesterdayStats, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	ui.yesterdayDailyDesc = widget.NewLabelWithStyle(formatDailyTotal, fyne.TextAlignTrailing, fyne.TextStyle{})
	ui.yesterdayDailyLabel = widget.NewLabelWithStyle("0:00:00", fyne.TextAlignLeading, fyne.TextStyle{})
	ui.yesterdayStartDesc = widget.NewLabelWithStyle(formatFirstStart, fyne.TextAlignTrailing, fyne.TextStyle{})
	ui.yesterdayStartLabel = widget.NewLabelWithStyle("No data", fyne.TextAlignLeading, fyne.TextStyle{})

	// Create buttons
	ui.startButton = widget.NewButton(textStart, ui.handleStartStop)
	ui.breakButton = widget.NewButton(textStartBreak, ui.handleBreak)
	ui.cancelButton = widget.NewButton(textCancel, ui.handleCancel)
	ui.breakButton.Disable()
	ui.cancelButton.Disable()
}

func (ui *UI) layoutWidgets() {
	// Create grid for today's stats (2 columns)
	todayGrid := container.NewGridWithColumns(2)
	todayGrid.Add(ui.todayTimeDesc)
	todayGrid.Add(ui.todayTimeLabel)
	todayGrid.Add(ui.breakDesc)
	todayGrid.Add(ui.breakLabel)
	todayGrid.Add(ui.dailyDesc)
	todayGrid.Add(ui.dailyLabel)
	todayGrid.Add(ui.firstStartDesc)
	todayGrid.Add(ui.firstStartLabel)

	// Create grid for yesterday's stats (2 columns)
	yesterdayGrid := container.NewGridWithColumns(2)
	yesterdayGrid.Add(ui.yesterdayDailyDesc)
	yesterdayGrid.Add(ui.yesterdayDailyLabel)
	yesterdayGrid.Add(ui.yesterdayStartDesc)
	yesterdayGrid.Add(ui.yesterdayStartLabel)

	// Create grid for both columns (2 columns)
	mainGrid := container.NewGridWithColumns(2)
	mainGrid.Add(todayGrid)
	mainGrid.Add(container.NewVBox(
		ui.yesterdayTitle,
		yesterdayGrid,
	))

	// Create grid for weekly total (2 columns, spans full width)
	weeklyGrid := container.NewGridWithColumns(2)
	weeklyGrid.Add(ui.weeklyDesc)
	weeklyGrid.Add(ui.weeklyLabel)

	// Create button container with vertical layout
	buttons := container.NewVBox(
		ui.startButton,
		ui.breakButton,
		ui.cancelButton,
	)

	// Layout everything vertically
	content := container.NewVBox(
		mainGrid,
		weeklyGrid,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			buttons,
			layout.NewSpacer(),
		),
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
			ui.todayTimeLabel.SetText(ui.formatDuration(ui.timer.GetTodaySessionTime()))
			ui.breakLabel.SetText(ui.formatDuration(ui.timer.GetCurrentBreakTime()))
			ui.weeklyLabel.SetText(ui.formatDuration(ui.timer.GetWeeklyTime()))
			ui.dailyLabel.SetText(ui.formatDuration(ui.timer.GetDailyTime()))
			ui.firstStartLabel.SetText(ui.timer.GetDayFirstStartTime())
			ui.yesterdayDailyLabel.SetText(ui.formatDuration(ui.timer.YesterdayTotal))
			ui.yesterdayStartLabel.SetText(ui.timer.GetYesterdayFirstStartTime())

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
