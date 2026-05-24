# logslice

Fast log filtering and time-range extraction tool for large structured log files.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

---

## Usage

```bash
# Extract logs between two timestamps
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" app.log

# Filter by log level
logslice --level error --from "2024-01-15T08:00:00Z" app.log

# Read from stdin and filter by keyword
cat app.log | logslice --grep "timeout" --from "2024-01-15T08:00:00Z"

# Output to a file
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" app.log -o out.log
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--level` | Filter by log level (info, warn, error) |
| `--grep` | Filter lines containing a string |
| `-o` | Write output to a file |

---

## Features

- Handles large log files efficiently with streaming reads
- Supports JSON and common structured log formats
- Binary search on sorted log files for fast range extraction
- Pipeable — works seamlessly with `grep`, `jq`, and other Unix tools

---

## License

MIT © 2024 yourusername