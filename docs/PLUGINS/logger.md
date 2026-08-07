# Logger Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The Logger plugin routes Go log messages to the native platform logger with configurable tag and severity levels.

---

## Quick Example

```go
package main

import (
    "github.com/sweet-juice/sweetjuice/plugins/logger"
)

func init() {
    log := logger.NewPlugin("myapp")
    log.Init(app)
}

func doWork() {
    logger.Debug("Starting operation %s", "fetch")
    logger.Info("User %s logged in", "alice")
    logger.Warn("Disk space low: %d MB", 100)
    logger.Error("Failed to connect: %v", err)
}
```

---

## API Reference

### `logger.NewPlugin(tag string) *LoggerPlugin`

Creates a new logger plugin instance.

| Argument | Type | Description |
|----------|------|-------------|
| `tag` | `string` | Tag prefix for log messages |

---

### `(*LoggerPlugin).Debug(format string, a ...interface{})`

Logs a debug message.

---

### `(*LoggerPlugin).Info(format string, a ...interface{})`

Logs an info message.

---

### `(*LoggerPlugin).Warn(format string, a ...interface{})`

Logs a warning message.

---

### `(*LoggerPlugin).Error(format string, a ...interface{})`

Logs an error message.

---

## Notes

- On Android, logs go to `Logcat` with the specified tag.
- On iOS, logs go to the system console.
- No special permissions required.
