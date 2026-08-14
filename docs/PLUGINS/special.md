# Special Permissions Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The Special Permissions plugin handles Android permissions that cannot be granted via normal runtime dialogs. These require opening system settings screens where the user manually enables them.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/special"
)

func requestAccessibility() {
    special.RequestAccessibility()
}

func openAppSettings() {
    special.RequestAppSettings()
}

func checkAllFiles() {
    granted, _ := special.CheckAllFilesAccess()
    fmt.Println("All files access:", granted)
}
```

---

## API Reference

### `special.NewPlugin() *SpecialPlugin`

Creates a new special permissions plugin instance.

---

### `(*SpecialPlugin).Request(perm SpecialType) (string, error)`

Opens the system settings screen for the given special permission type.

| Argument | Type | Description |
|----------|------|-------------|
| `perm` | `SpecialType` | One of: `Accessibility`, `AllFilesAccess` |

**Returns:** `{"status":"launched"}` or error

---

### `(*SpecialPlugin).Check(perm SpecialType) (bool, error)`

Checks whether the given special permission is currently granted.

| Argument | Type | Description |
|----------|------|-------------|
| `perm` | `SpecialType` | One of: `Accessibility`, `AllFilesAccess` |

**Returns:** `(granted bool, err error)`

---

## Convenience Helpers

| Function | Description |
|----------|-------------|
| `RequestAccessibility() (string, error)` | Opens Accessibility settings |
| `RequestAllFilesAccess() (string, error)` | Opens All Files Access settings (Android 11+) |
| `RequestBatteryExemption() (string, error)` | Opens battery optimization exemption dialog |
| `RequestNotificationAccess() (string, error)` | Opens notification access settings |
| `RequestAppSettings() (string, error)` | Opens app-specific system settings |
| `CheckAccessibility() (bool, error)` | Checks if Accessibility is granted |
| `CheckAllFilesAccess() (bool, error)` | Checks if All Files Access is granted |
| `CheckBatteryExemption() (bool, error)` | Checks if battery exemption is granted |
| `CheckNotificationAccess() (bool, error)` | Checks if notification access is granted |

---

## Special Permission Types

| Type | Settings Screen | Android Version |
|------|----------------|-----------------|
| `Accessibility` | `Settings.ACTION_ACCESSIBILITY_SETTINGS` | All |
| `AllFilesAccess` | `Settings.ACTION_MANAGE_APP_ALL_FILES_ACCESS_PERMISSION` | 11+ |
| `BatteryExemption` | `Settings.ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS` | 6.0+ |
| `NotificationAccess` | `Settings.ACTION_NOTIFICATION_LISTENER_SETTINGS` | All |

---

## Android Notes

- These permissions **cannot** be requested via `requestPermissions()`.
- The user must manually enable them in the system settings screen.
- `AllFilesAccess` requires Android 11 (API 30) or higher.
- `NotificationAccess` requires a `NotificationListenerService` to be declared in the app manifest with the `BIND_NOTIFICATION_LISTENER_SERVICE` permission.
