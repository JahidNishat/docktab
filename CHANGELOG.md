# Changelog

All notable changes to this project will be documented in this file.

The format is inspired by [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Automated tag-based releases
- JSON output support
- Live watch mode
- Expanded filtering options

## [0.1.0] - 2026-05-16

### Added
- Executable `docktab` CLI entrypoint
- Root CLI command and version command
- Container, image, volume, and network table commands
- Runtime configuration through `DOCKTAB_*` environment variables
- Structured debug logging support
- Safer Docker metadata parsing helpers
- Unit tests for Docker helpers, config, table formatting, and command behavior
- CI checks for formatting, vetting, race-enabled tests, and builds