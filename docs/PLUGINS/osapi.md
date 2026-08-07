# OS API Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The OS API plugin retrieves operating system and device information, returning platform-specific fields.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/osapi"
)

func logDeviceInfo() {
    info, err := osapi.GetInfo()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("OS:", info.OS)
    fmt.Println("Version:", info.Version)
    fmt.Println("Device:", info.Device)
    fmt.Println("Manufacturer:", info.Manufacturer)
}
```

---

## API Reference

### `osapi.NewPlugin() *OsApiPlugin`

Creates a new OS API plugin instance.

---

### `(*OsApiPlugin).GetInfo() (OsInfo, error)`

Returns system and device information.

**Returns:**

```go
type OsInfo struct {
    OS           string // "android" or "ios"
    Version      string // OS version string
    SDK          int    // SDK/API level (Android) or system version (iOS)
    Release      string // OS release codename
    Codename     string // OS codename
    Manufacturer string // Device manufacturer
    Model        string // Device model
    Device       string // Device identifier
    IsPhysical   bool   // Physical device or simulator
}
```

---

## Notes

- Field availability varies by platform.
- Android returns: `sdk_int`, `release`, `codename`, `manufacturer`, `model`, `device`.
- iOS returns: `name`, `localized_model`, `identifier_for_vendor`, `is_physical_device`.
- No special permissions required.
