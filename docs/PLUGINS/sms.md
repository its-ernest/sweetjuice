# SMS Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The SMS plugin provides access to device SMS messages across inbox, sent, draft, and all folders.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/sms"
)

func readInbox() {
    messages, err := sms.GetInbox()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, msg := range messages {
        fmt.Printf("From: %s, Body: %s\n", msg.Address, msg.Body)
    }
}

func readAll() {
    all, err := sms.GetAll()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Total messages:", len(all))
}
```

---

## API Reference

### `sms.NewPlugin() *SmsPlugin`

Creates a new SMS plugin instance.

---

### `(*SmsPlugin).GetInbox() (SmsFolder, error)`

Retrieves inbox messages.

---

### `(*SmsPlugin).GetSent() (SmsFolder, error)`

Retrieves sent messages.

---

### `(*SmsPlugin).GetDrafts() (SmsFolder, error)`

Retrieves draft messages.

---

### `(*SmsPlugin).GetAll() (SmsFolder, error)`

Retrieves all messages across all folders.

---

## Message Fields

```go
type SmsMessage struct {
    ID        int64  `json:"id"`
    Address   string `json:"address"`
    Body      string `json:"body"`
    Timestamp int64  `json:"timestamp"` // Unix millis
    Read      bool   `json:"read"`
    Type      string `json:"type"` // inbox, sent, draft, etc.
}
```

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.READ_SMS" />
```

---

## Notes

- Android only.
- Requires `READ_SMS` runtime permission.
- Messages are returned in chronological order (newest first).
