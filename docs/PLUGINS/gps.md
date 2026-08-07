# GPS Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The GPS plugin provides access to device location data, including last known location and continuous monitoring.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/gps"
)

func getLocation() {
    loc, err := gps.GetCurrentLocation()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Lat: %f, Lng: %f\n", loc.Latitude, loc.Longitude)
}

func startTracking() {
    gps.StartMonitoring()
}

func stopTracking() {
    gps.StopMonitoring()
}
```

---

## API Reference

### `gps.NewPlugin() *GpsPlugin`

Creates a new GPS plugin instance.

---

### `(*GpsPlugin).GetCurrentLocation() (Location, error)`

Returns the last known location.

**Returns:**

```go
type Location struct {
    Latitude  float64
    Longitude float64
    Accuracy  float64
    Timestamp string
}
```

---

### `(*GpsPlugin).StartMonitoring() (string, error)`

Starts continuous location updates.

---

### `(*GpsPlugin).StopMonitoring() (string, error)`

Stops continuous location updates.

---

## Events

### `gps:changed`

Fired when location updates during monitoring. Payload:

```json
{
  "latitude": 37.7749,
  "longitude": -122.4194,
  "accuracy": 10.0,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />
```

---

## iOS Requirements

Add to `Info.plist`:

```xml
<key>NSLocationWhenInUseUsageDescription</key>
<string>We need your location to show nearby places</string>
```

---

## Notes

- `GetCurrentLocation()` returns the last known location, which may be stale.
- Use `StartMonitoring()` for real-time updates.
- Background location requires additional configuration and permissions.
