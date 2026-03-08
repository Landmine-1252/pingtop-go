# pingtop

[![CI](https://img.shields.io/github/actions/workflow/status/landmine-1252/pingtop-go/ci.yml?branch=main&label=ci&logo=githubactions)](https://github.com/Landmine-1252/pingtop-go/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/landmine-1252/pingtop-go?logo=go)](https://github.com/Landmine-1252/pingtop-go/blob/main/go.mod)
[![Release](https://img.shields.io/github/v/release/landmine-1252/pingtop-go?display_name=tag&sort=semver&logo=github)](https://github.com/Landmine-1252/pingtop-go/releases)
[![Coverage](https://img.shields.io/codecov/c/github/landmine-1252/pingtop-go?logo=codecov)](https://app.codecov.io/gh/landmine-1252/pingtop-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/landmine-1252/pingtop-go)](https://goreportcard.com/report/github.com/landmine-1252/pingtop-go)
[![Downloads](https://img.shields.io/github/downloads/landmine-1252/pingtop-go/latest/total?logo=github)](https://github.com/Landmine-1252/pingtop-go/releases)
[![Platforms](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-1f2937)](https://github.com/Landmine-1252/pingtop-go/releases)
[![Arch](https://img.shields.io/badge/arch-amd64%20%7C%20arm64%20%7C%20armv7-2563eb)](https://github.com/Landmine-1252/pingtop-go/releases)

`pingtop` is a Go rewrite of the original Python [`pingtop`](https://github.com/Landmine-1252/pingtop). The Go code, releases, and CI live in [`Landmine-1252/pingtop-go`](https://github.com/Landmine-1252/pingtop-go).

![pingtop demo](demo.gif)

## What It Does

- live terminal UI with redraws, color, and keyboard controls
- concurrent ping and DNS checks with simple failure classification
- interactive mode and headless mode
- single-binary releases for Linux, macOS, and Windows
- ad hoc target runs from the command line without writing CSV logs

## Downloads

Prebuilt binaries are published on [GitHub Releases](https://github.com/Landmine-1252/pingtop-go/releases).

- platforms: Linux, macOS, Windows
- architectures: `amd64`, `arm64`, Linux `armv7`
- packaging: single binary per archive, no Go runtime required
- size: release asset sizes vary by platform and version; see the latest release assets for exact download sizes

## Quick Start

Run from source:

```bash
go run .           # interactive UI when a supported TTY is available
go run . -n        # headless mode
go run . -o        # single pass
go run . -v        # version
go run . -h        # help
go run . 1.1.1.1
go run . example.com 1.1.1.1
```

Build locally:

```bash
go build -o pingtop
./pingtop
./pingtop -h
./pingtop -v
```

Build `pingtop.exe` on Windows:

```powershell
go build -o pingtop.exe .
.\pingtop.exe
.\pingtop.exe -h
.\pingtop.exe -v
```

Passing one or more positional targets overrides the configured target list for that run only. Those ad hoc runs keep the normal UI or headless behavior, but CSV logging is disabled for that session.

## Controls

Press `h` in the UI to show or hide the full help panel. Common controls:

- `q` or `Esc`: quit
- `h`: show or hide help
- `p`: pause or resume
- `a` / `d`: add or delete a target
- `s`: save a snapshot
- `u`: open the release page

## Release Flow

Before tagging a release, update [`internal/pingtop/version.go`](internal/pingtop/version.go) so `Version` matches the tag value without the optional leading `v`.

Then push either tag format:

```bash
git tag 0.1.3
git push origin 0.1.3
```

or:

```bash
git tag v0.1.3
git push origin v0.1.3
```

The release workflow verifies that the tag and source version match, runs tests, builds release archives, and attaches them to the GitHub release automatically.

## Tests

```bash
go test ./...
```
