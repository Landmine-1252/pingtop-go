# pingtop

`pingtop` is a Go rewrite of the original Python `pingtop` monitor. It keeps the same practical operator workflow:

- classify intermittent issues into general network, DNS, or isolated path failures
- ping configured IP and hostname targets on a schedule
- resolve hostnames separately so DNS failures stay distinct from reachability failures
- run in an interactive ANSI UI or in headless mode
- persist config, CSV logs, rotated logs, and text snapshot reports beside the launched binary

## Go-Specific Improvements

- single native binary, no Python runtime required
- goroutine-based concurrent checks
- `go run` uses the current working directory for runtime files instead of a temporary build directory
- Linux/WSL and Windows interactive mode use native terminal APIs without extra third-party dependencies

## Runtime Files

The Go version keeps the Python-compatible filenames:

- `pingtop.json`
- `pingtop_log.csv`
- `pingtop_snapshot_YYYYMMDD_HHMMSS.txt`

The config schema is intentionally compatible with the Python project.

## Run From Source

From the project folder:

```bash
go run . 
go run . --once
go run . --no-ui
```

## Build Locally

Build a binary for your current platform:

```bash
go build -o pingtop
./pingtop
./pingtop --once
./pingtop --no-ui
```

Build a Windows executable named `pingtop.exe`:

```powershell
go build -o pingtop.exe .
.\pingtop.exe
.\pingtop.exe --once
.\pingtop.exe --no-ui
```

Prebuilt binaries can also be published from GitHub Releases.

## GitHub Releases

When you push a tag in `vX.Y.Z` format, the release workflow builds and publishes archives for:

- Linux `amd64`
- Linux `arm64`
- Linux `armv7`
- macOS `amd64`
- macOS `arm64`
- Windows `amd64`
- Windows `arm64`

Release binaries embed the tag version and the GitHub repository URL so the in-app update check points at the public repo automatically.

## Controls

- `q`: quit
- `p`: pause/resume
- `l`: cycle logging mode
- `+` or `=` / `-` or `_`: adjust check interval
- `>` or `.` / `<` or `,`: adjust UI refresh interval
- `a`: add a target
- `d`: delete a target
- `w`: change around-failure window
- `t`: change rolling stats window
- `r`: reset counters
- `s`: save a snapshot report
- `u`: open the project release page
- `h`: show/hide help

## Logging Modes

- `all`
- `failures_only`
- `around_failure`

Around-failure mode keeps a rolling pre-failure buffer and captures post-failure results. Logs rotate automatically when `log_rotation_max_mb` is exceeded.

## Update Checks

If `update_check_enabled` is true and `update_repo_url` points to a GitHub repo, `pingtop` checks the latest release metadata in the background and surfaces availability in the UI. It does not self-update.

## Tests

Run the Go test suite with:

```bash
go test ./...
```

## Current Platform Notes

- Linux/WSL interactive mode is implemented.
- Windows interactive mode is implemented, including ANSI redraw/color support in modern consoles.
- Other platforms fall back to headless mode if interactive input is unavailable.
