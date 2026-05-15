# DOCKTAB

> A fast, beautiful Docker CLI tool that displays clean, professional tables.

`docktab` turns messy Docker output into readable, colored tables with powerful filtering and sorting.

## Features

- `docktab ps` — Beautiful container table with filtering & sorting
- `docktab images` — Clean image listing
- `docktab volumes` — Volume management table
- `docktab networks` — Network listing
- Smart truncation and terminal width detection
- Status with color coding
- Human-readable sizes and timestamps

## Installation

```bash
go install github.com/JahidNishat/docktab/cmd/docktab@latest
```
Or build from source:
```bash
git clone https://github.com/JahidNishat/docktab.git
cd docktab
make build
```
## Usage Examples
```bash
# Containers
docktab ps
docktab ps -a --compact
docktab ps --name redis --sort created

# Images
docktab images
docktab images --sort size --full

# Volumes
docktab volumes
docktab volumes --sort driver

# Networks
docktab networks
```

## Architecture
- Clean, interface-driven design:
- internal/command — Command interface
- internal/registry — Explicit registration
- Easy to extend with new commands

## Roadmap

- [ ] --watch mode
- [ ] Interactive TUI mode (--interactive)
- [ ] Better configuration system