# Biometric Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The Biometric plugin provides biometric authentication support, allowing apps to check device capability and trigger fingerprint, face recognition, or other biometric prompts.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/biometric"
)

func authenticate() {
    canAuth, _ := biometric.CanAuthenticate()
    if canAuth.Can {
        res, err := biometric.Authenticate(biometric.AuthOptions{
            Title:       "Login Required",
            Description: "Please authenticate to continue",
        })
        fmt.Println("Auth result:", res)
    }
}
```

---

## API Reference

### `biometric.NewPlugin() *BiometricPlugin`

Creates a new biometric plugin instance.

---

### `(*BiometricPlugin).CanAuthenticate() (CanAuthResult, error)`

Checks whether biometric authentication is available on the device.

**Returns:** `CanAuthResult` with fields:
- `Can bool` — whether biometric auth is available
- `Error string` — reason if unavailable

---

### `(*BiometricPlugin).Authenticate(options AuthOptions) (string, error)`

Triggers a biometric authentication prompt.

| Field | Type | Description |
|-------|------|-------------|
| `Title` | `string` | Dialog title |
| `Description` | `string` | Dialog description |
| `NegativeButtonText` | `string` | Text for cancel button (optional) |

**Returns:** JSON string with authentication result

---

## Events

### `biometric:result`

Fired when authentication completes. Payload:

```json
{
  "success": true,
  "error": ""
}
```

Or on failure:

```json
{
  "success": false,
  "error": "USER_CANCELLED"
}
```

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.USE_BIOMETRIC" />
```

Ensure `androidx.biometric:biometric` is in your `app/build.gradle`.

---

## iOS Requirements

Add to `Info.plist`:

```xml
<key>NSFaceIDUsageDescription</key>
<string>Use Face ID to authenticate</string>
<key>NSBiometricUsageDescription</key>
<string>Use biometric to authenticate</string>
```
