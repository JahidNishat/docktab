# docktab

> A fast, beautiful Docker table viewer for developers who live in the terminal.

`docktab` turns messy Docker CLI output into clean, professional tables with powerful filtering and sorting.

## Features (v0.1)

- `docktab ps` — Beautiful container table
- Clean architecture with interfaces
- Easy command extensibility

## Architecture

This project demonstrates clean Go architecture:

- `internal/command` — Command interface
- `internal/registry` — Explicit registration
- Nested command packages (`commands/ps/`, etc.)

## Installation

```bash
git clone https://github.com/JahidNishat/docktab.git
cd docktab
make build
```

## Usage

```bash
docktab ps
docktab ps -a --compact
```
