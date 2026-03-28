# awf-plugin-time

AWF plugin that returns system date/time for workflow injection.

## Install

```bash
awf plugin install awf-project/awf-plugin-time
```

## Operation: `time`

Returns the current system date and time.

### Inputs

| Name | Type | Default | Description |
|------|------|---------|-------------|
| `format` | string | RFC3339 | Go [time layout](https://pkg.go.dev/time#pkg-constants) |
| `tz` | string | Local | IANA timezone name |

### Output

- **Output string:** Formatted date/time
- **Data:**
  - `timestamp` — RFC3339 (always, regardless of format)
  - `unix` — Unix timestamp (int64)
  - `timezone` — Resolved timezone name

### Workflow Example

```yaml
name: get-current-time
description: Fetch system time via plugin

states:
  get_time:
    type: operation
    operation: awf-plugin-time.time
    inputs:
      format: "2006-01-02 15:04:05"
      tz: "UTC"
    on_success: done

  done:
    type: terminal
    status: success
    message: "Current time: {{states.get_time.output}}"
```

## Development

```bash
make build      # Build binary
make test       # Run tests
make install    # Install to ~/.local/share/awf/plugins/time/
make clean      # Remove binary
```

## License

MIT
