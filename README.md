# Time Tracker

A desktop application for tracking working hours with pause functionality. The application features an always-on-top window that helps you keep track of your work time and breaks.

## Features

- Start/Stop work timer
- Track pause time separately
- Reset functionality
- Session state persistence
- Weekly time tracking (Monday-based)
- Always-on-top window

## Requirements

- Go 1.21 or later
- Fyne GUI toolkit

## Installation

1. Install Go from [https://go.dev/dl/](https://go.dev/dl/)
2. Clone this repository
3. Install dependencies:
```bash
go mod download
```
4. Build and run:
```bash
go run main.go
```

## Building

To build a standalone executable:

```bash
go build
```

## Usage

- Click the "Start/Stop" button to begin/end work time tracking
- Use the "Pause" button to track break time
- "Reset" button clears the current session
- The application automatically saves state on pause or reset
- Weekly statistics are automatically tracked and displayed 