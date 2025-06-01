# Contributing to Time Tracker

Thank you for your interest in contributing to Time Tracker! This document provides guidelines and instructions for contributing to the project.

## Development Setup

1. **Prerequisites**
   - Go 1.23.0 or later
   - Make (or mingw32-make for Windows)
   - GCC compiler (MinGW-w64 for Windows)
   - For Linux: `xorg-dev` and `libgl1-mesa-dev` packages

2. **Clone the Repository**
   ```bash
   git clone https://github.com/LitschiW/timetracker.git
   cd timetracker
   ```

3. **Install Dependencies**
   ```bash
   make deps
   ```

## Development Commands

The project includes a Makefile with several useful commands:

```bash
make build       # Build the application
make run        # Build and run the application
make test       # Run all tests
make clean      # Clean build artifacts
make deps       # Install/update dependencies
make screenshots # Generate application screenshots
```

On Windows, you can use either `make` or `mingw32-make` depending on your setup.

## Development Workflow

1. **Create a New Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**
   - Write your code
   - Add tests if applicable
   - Update documentation if needed

3. **Test Your Changes**
   ```bash
   make test
   ```

4. **Generate Screenshots** (if UI changes)
   ```bash
   make screenshots
   ```

5. **Build and Test Locally**
   ```bash
   make run
   ```

6. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "Description of your changes"
   ```

7. **Push and Create a Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```
   Then create a pull request on GitHub.

## Code Style Guidelines

- Follow standard Go formatting (use `gofmt`)
- Add comments for non-obvious code
- Keep functions focused and small
- Write meaningful commit messages

## Testing

- Add tests for new features
- Ensure existing tests pass
- UI tests should use the screenshot test framework

## Documentation

- Update README.md if adding new features
- Document new commands or configuration options
- Keep code comments up to date

## Need Help?

If you have questions or need help with setup, feel free to:
- Open an issue
- Ask questions in pull requests
- Reach out to maintainers

Thank you for contributing! 