# Time Tracker

A desktop application for tracking working hours with pause functionality. The application features an always-on-top window that helps you keep track of your work time and breaks.

## Screenshots

Here's how the application looks in action:

| ⏸️ Initial State | ⏳ Running Working Session | ☕ Break Tracking |
|:---:|:---:|:---:|
| ![Initial State](.github/screenshots/initial_state.png?raw=true) | ![Running Working Session](.github/screenshots/working_session.png?raw=true) | ![Break Tracking](.github/screenshots/break_state.png?raw=true) |

The application maintains a clean, focused interface that helps you track your work time efficiently.

## Features

- Start/Stop work timer
- Track pause time separately
- Reset functionality
- Session state persistence
- Weekly time tracking (Monday-based)
- Always-on-top window

## Requirements

- Go 1.23.0 or later

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

# For Linux/macOS
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
- "Reset" button stops the current session
- The application automatically saves state on pause or reset
- Weekly statistics are automatically tracked and displayed 

## Contributing

This project was primarily "vibe coded" - built with a focus on getting things working and iterating quickly. While this approach helped us move fast and ship features, there's always room for improvement! We welcome pull requests to:

- Clean up and refactor the code
- Add new features
- Improve documentation
- Fix bugs
- Enhance the UI/UX

Check out our [Contributing Guide](.github/CONTRIBUTING.md) to get started! 