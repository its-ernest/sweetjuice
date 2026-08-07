# Device State Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The Device State plugin provides access to device state including battery level, charging status, and connectivity information. Supports optional continuous monitoring.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/devicestate"
)

func checkState() {
    state, err := devicestate.GetState()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Battery:", state.BatteryLevel, "%")
    fmt.Println("Charging:", state.IsCharging)
    fmt.Println("Connected:", state.IsConnected)
}

func startMonitoring() {
    devicestate.StartMonitoring()
}

func stopMonitoring() {
    devicestate.StopMonitoring()
}
```

---

## API Reference

### `devicestate.NewPlugin() *DeviceStatePlugin`

Creates a new device state plugin instance.

---

### `(*DeviceStatePlugin).GetState() (DeviceState, error)`

Returns the current device state.

**Returns:**

```go
type DeviceState struct {
    BatteryLevel int    // 0-100
    IsCharging   bool
    IsConnected  bool   // network connectivity
    Connectivity string // wifi, cellular, none
}
```

---

### `(*DeviceStatePlugin).StartMonitoring() (string, error)`

Starts continuous monitoring of device state changes.

---

### `(*DeviceStatePlugin).StopMonitoring() (string, error)`

Stops continuous monitoring.

---

## Events

### `devicestate:changed`

Fired when device state changes during monitoring. Payload:

```json
{
  "battery_level": 85,
  "is_charging": true,
  "is_connected": true,
  "connectivity": "wifi"
}
```

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
```

---

## Notes

- Android only.
- `GetState()` returns a snapshot. Use `StartMonitoring()` for real-time updates.
