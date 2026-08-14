# WorkManager Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The WorkManager plugin manages background tasks via the Android WorkManager API, supporting constraint-based one-time and periodic scheduling.

---

## Quick Example

```go
package main

import (
    "github.com/sweet-juice/sweetjuice/plugins/workmanager"
)

func registerTask() {
    workmanager.RegisterTask("cleanup", func(ctx workmanager.TaskContext) error {
        fmt.Println("Running cleanup task")
        return nil
    })
}

func scheduleCleanup() {
    workmanager.EnqueueOneTime("cleanup", workmanager.DefaultConstraints())
}

func schedulePeriodic() {
    workmanager.EnqueuePeriodic("cleanup", 15, workmanager.DefaultConstraints(), false)
}
```

---

## API Reference

### `workmanager.NewPlugin() *WorkManagerPlugin`

Creates a new WorkManager plugin instance.

---

### `(*WorkManagerPlugin).RegisterTask(name string, fn TaskFunc)`

Registers a task function that can be executed by WorkManager.

```go
type TaskFunc func(ctx TaskContext) error
```

---

### `(*WorkManagerPlugin).EnqueueOneTime(taskKey string, constraints *Constraints) (string, error)`

Schedules a one-time background task.

| Argument | Type | Description |
|----------|------|-------------|
| `taskKey` | `string` | Unique task identifier |
| `constraints` | `*Constraints` | Execution constraints |

---

### `(*WorkManagerPlugin).EnqueuePeriodic(taskKey string, intervalMinutes int, constraints *Constraints, replaceExisting bool) (string, error)`

Schedules a repeating background task using `enqueueUniquePeriodicWork`.

| Argument | Type | Description |
|----------|------|-------------|
| `taskKey` | `string` | Unique task identifier used as the WorkManager unique name |
| `intervalMinutes` | `int` | Repeat interval in minutes (minimum 15) |
| `constraints` | `*Constraints` | Execution constraints |
| `replaceExisting` | `bool` | If `true`, replaces any existing periodic work with the same `taskKey` (`REPLACE`). If `false`, keeps the existing work (`KEEP`) |

---

### `(*WorkManagerPlugin).IsEnqueued(taskKey string) (bool, error)`

Checks if a task is currently scheduled.

---

### `(*WorkManagerPlugin).CancelAll() string`

Cancels all scheduled tasks.

---

### `workmanager.DefaultConstraints() Constraints`

Returns default constraints (requires network, battery not low).

---

## Task Context

```go
type TaskContext struct {
    Data     map[string]string // Input data from enqueue
    SetResult(result string)    // Set task result
}
```

---

## Events

### `workmanager:execute`

Fired when a background task executes. Payload includes task key and input data.

### `workmanager:result`

Fired when a task completes. Payload:

```json
{
  "task_key": "cleanup",
  "result": "success",
  "error": ""
}
```

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_SPECIAL_USE" />
```

---

## Notes

- Android only. iOS uses BGTaskScheduler instead.
- Tasks must be registered before they can be enqueued.
- Tasks run even if the app is killed (system-managed).
