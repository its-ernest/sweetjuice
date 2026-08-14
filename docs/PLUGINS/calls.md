# Calls Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The Calls plugin provides access to the device call log, allowing retrieval of recent or all call records.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/calls"
)

func loadRecentCalls() {
    calls, err := calls.GetRecent(10)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, call := range calls {
        fmt.Printf("%s: %s (%s)\n", call.Number, call.Type, call.Date)
    }
}

func loadAllCalls() {
    calls, err := calls.GetAll()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Total calls:", len(calls))
}
```

---

## API Reference

### `calls.NewPlugin() *CallsPlugin`

Creates a new calls plugin instance.

---

### `(*CallsPlugin).GetRecent(limit int) (CallLog, error)`

Retrieves the most recent call records.

| Argument | Type | Description |
|----------|------|-------------|
| `limit` | `int` | Maximum number of records to return |

**Returns:** `CallLog` (slice of call records) or error

---

### `(*CallsPlugin).GetLast(limit int) (CallLog, error)`

An alias for `GetRecent(limit)`.

---

### `(*CallsPlugin).GetAll() (CallLog, error)`

Retrieves all call records.

**Returns:** `CallLog` or error

---

## Call Record Fields

```go
type CallRecord struct {
    ID        int64  `json:"id"`
    Number    string `json:"number"`
    Type      string `json:"type"`      // INCOMING, OUTGOING, MISSED
    Date      string `json:"date"`      // timestamp
    Duration  int    `json:"duration"`  // seconds
    CachedName string `json:"cached_name,omitempty"`
    Geo       string `json:"geo,omitempty"`
}
```

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.READ_CALL_LOG" />
```

---

## Notes

- Android only.
- Requires `READ_CALL_LOG` runtime permission.
- Call types are returned as strings: `INCOMING`, `OUTGOING`, `MISSED`.
