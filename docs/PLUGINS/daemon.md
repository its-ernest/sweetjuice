# Daemon Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The Daemon plugin manages the application lifecycle as a background/foreground service. On Android Oreo+, it starts a Foreground Service with a notification to keep the app alive.

---

## Quick Example

```go
package main

import (
    "github.com/sweet-juice/sweetjuice/plugins/daemon"
)

func startBackground() {
    _, err := daemon.NewPlugin().Start(daemon.Options{
        NotificationTitle: "Sweet Juice",
        NotificationText:  "Running in background",
    })
    if err != nil {
        fmt.Println("Error:", err)
    }
}

func stopBackground() {
    daemon.NewPlugin().Stop()
}
```

---

## API Reference

### `daemon.NewPlugin() *DaemonPlugin`

Creates a new daemon plugin instance.

---

### `(*DaemonPlugin).Start(opts Options) (string, error)`

Starts the background service.

| Field | Type | Description |
|-------|------|-------------|
| `NotificationTitle` | `string` | Foreground notification title |
| `NotificationText` | `string` | Foreground notification text |
| `ChannelID` | `string` | Notification channel ID (optional) |
| `ChannelName` | `string` | Notification channel name (optional) |

---

### `(*DaemonPlugin).Stop() (string, error)`

Stops the background service.

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_SPECIAL_USE" />
```

Service must be declared in `AndroidManifest.xml`:

```xml
<service
    android:name="com.sweetjuice.pkg.daemon.SweetJuiceDaemonService"
    android:enabled="true"
    android:exported="false"
    android:foregroundServiceType="specialUse">
    <property android:name="android.app.PROPERTY_SPECIAL_USE_FGS_SUBTYPE"
        android:value="SweetJuice Background Logic" />
</service>
```

---

## Notes

- Android only. iOS handles background differently via BGTaskScheduler.
- On Android Oreo+, a foreground service is mandatory for persistent background logic.
