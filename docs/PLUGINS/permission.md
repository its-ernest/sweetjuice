# Permission Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The Permission plugin provides runtime permission management for Android. It allows your Go code to check permission status and request one or more runtime permissions. Results are delivered asynchronously via the `permissions:changed` event.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/permission"
)

func requestCamera() {
    status, _ := permission.Check("android.permission.CAMERA")
    if status != "granted" {
        permission.Request("android.permission.CAMERA")
    }
}

func requestMultiple() {
    permission.RequestMultiple([]string{
        "android.permission.READ_CALL_LOG",
        "android.permission.READ_SMS",
    })
}
```

---

## API Reference

### `permission.NewPlugin() *PermissionPlugin`

Creates a new permission plugin instance.

---

### `(*PermissionPlugin).Init(app *core.Application) error`

Initializes the plugin with the application context. Registers the `permissions:result` native callback.

---

### `(*PermissionPlugin).Check(permission string) (string, error)`

Checks whether a single permission is currently granted.

| Argument | Type | Description |
|----------|------|-------------|
| `permission` | `string` | Android permission constant, e.g. `android.permission.CAMERA` |

**Returns:** `"granted"` or `"denied"`

---

### `(*PermissionPlugin).Request(permission string) (string, error)`

Requests a single runtime permission. Triggers the native system dialog.

| Argument | Type | Description |
|----------|------|-------------|
| `permission` | `string` | Android permission constant |

**Returns:** `{"status":"requested"}` or error

---

### `(*PermissionPlugin).RequestMultiple(permissions []string) (string, error)`

Requests multiple runtime permissions in one dialog.

| Argument | Type | Description |
|----------|------|-------------|
| `permissions` | `[]string` | Array of Android permission constants |

**Returns:** `{"status":"requested"}` or error

---

## Events

### `permissions:changed`

Fired when a permission request completes. Payload:

```json
{
  "permission": "android.permission.CAMERA",
  "granted": true
}
```

---

## Android Manifest

Declare all permissions you intend to request:

```xml
<uses-permission android:name="android.permission.CAMERA" />
<uses-permission android:name="android.permission.READ_SMS" />
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
```

---

## Notes

- Android only. iOS does not use runtime permissions in the same way.
- Background-restricted permissions like `ACCESS_BACKGROUND_LOCATION` must be requested separately after foreground permissions are granted.
- The `permissions:changed` event fires once per permission in a `RequestMultiple` call.
