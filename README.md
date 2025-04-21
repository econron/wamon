# Wamon 🦭

A CLI tool where a ringed seal (ワモンアザラシ) praises you for tracking your daily activities!

![Go](https://img.shields.io/badge/Go-1.22-blue)
![License: MIT](https://img.shields.io/badge/license-MIT-green)
![PRs welcome](https://img.shields.io/badge/PRs-welcome-brightgreen)

## Overview

Wamon helps you track your daily research and programming activities through a friendly CLI interface. The app features a cute ringed seal mascot that provides encouragement and keeps track of your work.

## Features

- Interactive CLI interface
- Track research activities
- Record programming accomplishments
- List previous entries
- Filter records by category
- Satisfaction rating system
- Encouraging seal messages

## Installation

### Using Go (recommended)

```bash
go install github.com/econron/wamon@latest
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/econron/wamon.git
cd wamon
```

2. Build the executable:
```bash
go build -o wamon
```

3. Move the executable to your PATH:
```bash
sudo mv wamon /usr/local/bin/
```

## Usage

### Interactive Mode

Simply run the `wamon` command without arguments to enter interactive mode:

```bash
wamon
```

This will start an interactive session where you can:
1. Choose a category (Research, Programming, or both)
2. Enter details about your activity
3. Rate your satisfaction
4. Receive encouragement from the seal!

### Listing Previous Entries

To list all your previous entries:

```bash
wamon list
```

Filter by category:

```bash
wamon list -c "調べ物"  # or you can input "research"
wamon list -c "プログラマ"  # or you can input programming
wamon list -c "調べてプログラマ"  # Both
```

## Configuration

By default, Wamon stores data in `~/.wamon/wamon.db`. You can customize this location:

```bash
wamon --db /custom/path/to/database.db
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

