# Contributing to Time Tracker

Thank you for your interest in contributing to Time Tracker! We welcome contributions from everyone.

## How to Contribute

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Make your changes
4. Write or update tests as needed
5. Run the test suite to ensure everything passes
6. Submit a pull request

## Development Setup

1. Install Go 1.23.0 or later
2. Install required dependencies:
   ```bash
   # For Linux
   sudo apt-get update
   sudo apt-get install -y xorg-dev libgl1-mesa-dev
   ```
3. Clone your fork
4. Run `go mod download` to install dependencies
5. Run `go test ./...` to ensure everything works

## Code Structure

- `src/` - Main application source code
- `tests/` - Test files
- `.github/` - GitHub specific files (workflows, templates, etc.)

## Pull Request Guidelines

1. Keep changes focused and atomic
2. Follow existing code style
3. Include tests for new features
4. Update documentation as needed
5. Add a clear description of your changes

## Need Help?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the codebase

We appreciate all contributions, whether it's:
- Code improvements
- Documentation updates
- Bug fixes
- Feature additions
- UI/UX enhancements 