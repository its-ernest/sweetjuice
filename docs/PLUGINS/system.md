# System Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The System plugin exposes OS-level information, replacing the older `osapi` plugin with a clearer API surface. It currently provides device and OS metadata through a single `GetInfo()` call.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/system"
)

func logDeviceInfo() {
    info, err := system.NewPlugin().GetInfo()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("System:", info.SystemName)
    fmt.Println("Version:", info.SystemVersion)
    fmt.Println("Model:", info.Model)
    fmt.Println("Manufacturer:", info.Manufacturer)
}
```

---

## API Reference

### `system.NewPlugin() *SystemPlugin`

Creates a new System plugin instance.

---

### `(*SystemPlugin).GetInfo() (SystemInfo, error)`

Returns OS and device information from the native layer.

**Returns:** `SystemInfo` struct or error.

---

### `SystemInfo`

```go
type SystemInfo struct {
    SystemName          string `json:"system_name"`
    SystemVersion       string `json:"system_version"`
    Model               string `json:"model"`
    SdkInt              int    `json:"sdk_int"`
    Release             string `json:"release"`
    Codename            string `json:"codename"`
    Manufacturer        string `json:"manufacturer"`
    Brand               string `json:"brand"`
    Board               string `json:"board"`
    Device              string `json:"device"`
    Product             string `json:"product"`
    Hardware            string `json:"hardware"`
    BaseOS              string `json:"base_os"`
    SecurityPatch       string `json:"security_patch"`
    Name                string `json:"name"`
    LocalizedModel      string `json:"localized_model"`
    IdentifierForVendor string `json:"identifier_for_vendor"`
    IsPhysicalDevice    bool   `json:"is_physical_device"`
}
```

| Field | Platform | Description |
|-------|----------|-------------|
| `system_name` | All | `"android"` or `"ios"` |
| `system_version` | All | Human-readable OS version |
| `model` | All | Device model name |
| `sdk_int` | Android | `Build.VERSION.SDK_INT` |
| `release` | Android | `Build.VERSION.RELEASE` |
| `codename` | Android | `Build.VERSION.CODENAME` |
| `manufacturer` | Android | `Build.MANUFACTURER` |
| `brand` | Android | `Build.BRAND` |
| `board` | Android | `Build.BOARD` |
| `device` | Android | `Build.DEVICE` |
| `product` | Android | `Build.PRODUCT` |
| `hardware` | Android | `Build.HARDWARE` |
| `base_os` | Android | `Build.VERSION.BASE_OS` |
| `security_patch` | Android | `Build.VERSION.SECURITY_PATCH` |
| `name` | iOS | Device name |
| `localized_model` | iOS | Localized model string |
| `identifier_for_vendor` | iOS | Vendor identifier |
| `is_physical_device` | iOS | Physical device or simulator |

---

## Android Requirements

No special permissions required for basic system info.

---

## Notes

- Replaces the deprecated `osapi` plugin.
- Field availability varies by platform; unpopulated fields will be zero-value.
- iOS support requires the matching Swift plugin implementation in the Xcode project.
