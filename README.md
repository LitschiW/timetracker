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

### Pre-built Binaries

You can download pre-built binaries for Windows and Linux from the [Releases](../../releases) page.

### Building from Source

1. Install Go from [https://go.dev/dl/](https://go.dev/dl/)
2. Clone this repository
3. Install dependencies:
```bash
go mod download
```
4. Build and run:
```bash
go run .
```

## Building Standalone Executable

To build a standalone executable that can be distributed and run without Go installed:

```bash
# For Windows
go build -o timetracker.exe

# For Linux
go build -o timetracker

# For macOS
go build -o timetracker
```

The resulting executable:
- Can be moved to any location
- Doesn't require Go to be installed
- Will create its data file (`timetracker.json`) in the directory it's run from
- Works as a portable application

### Linux Dependencies

If building or running on Linux, you'll need to install some additional dependencies:
```bash
sudo apt-get update
sudo apt-get install -y xorg-dev libgl1-mesa-dev
```

## Continuous Integration

This project uses GitHub Actions to automatically build executables for Windows and Linux. The builds are triggered on:
- Every push to the main branch
- Every pull request to the main branch
- When creating a new release

The built executables are available:
- As artifacts from each workflow run
- As downloadable assets on the releases page when a new release is created

## Usage

- Click the "Start/Stop" button to begin/end work time tracking
- Use the "Pause" button to track break time
- "Reset" button clears the current session
- The application automatically saves state on pause or reset
- Weekly statistics are automatically tracked and displayed 