# cronview

Terminal dashboard for monitoring cron job history and failure rates on Linux servers.

![Go Version](https://img.shields.io/badge/go-1.21+-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Installation

```bash
go install github.com/yourusername/cronview@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/cronview/releases).

## Usage

Run `cronview` in your terminal to launch the interactive dashboard:

```bash
cronview
```

By default, cronview reads from `/var/log/syslog` and `/var/log/cron`. You can specify a custom log path:

```bash
cronview --log /var/log/cron.log
```

Filter by a specific time range:

```bash
cronview --since 24h
```

### Dashboard Controls

| Key | Action |
|-----|--------|
| `↑ / ↓` | Navigate job list |
| `Enter` | View job details and run history |
| `f` | Filter by failure status |
| `r` | Refresh data |
| `q` | Quit |

## Requirements

- Linux (systemd or syslog-based cron logging)
- Go 1.21+ (if building from source)

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

This project is licensed under the [MIT License](LICENSE).