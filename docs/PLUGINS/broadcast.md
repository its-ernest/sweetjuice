# Broadcast Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The Broadcast plugin bridges Android's Intent-based broadcast system to Go. It allows you to listen for system events (like battery status or boot completion) and send system-wide signals.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/broadcast"
)

func init() {
    // Listen for power connected
    broadcast.On("android.intent.action.ACTION_POWER_CONNECTED", func(data interface{}) {
        fmt.Println("Phone plugged in!")
    })
}

func triggerAction() {
    // Send a custom intent to the system
    broadcast.Post("com.myapp.CUSTOM_ACTION", map[string]interface{}{
        "timestamp": 123456789,
        "status":    "active",
    })
}
```

---

## API Reference

### `broadcast.NewPlugin() *BroadcastPlugin`

Creates a new Broadcast plugin instance.

---

### `(*BroadcastPlugin).On(action string, callback func(interface{}))`

Registers a listener for a specific Android Intent action.

- **System Actions**: For system-wide broadcasts (like `BOOT_COMPLETED`), ensure the action is registered in your `AndroidManifest.xml`.
- **Dynamic Actions**: For common runtime actions, the plugin registers a dynamic `BroadcastReceiver` automatically.

---

### `(*BroadcastPlugin).Post(action string, extras map[string]interface{})`

Sends a system-wide broadcast (Intent) with optional extras.

| Argument | Type | Description |
|----------|------|-------------|
| `action` | `string` | The Intent action string |
| `extras` | `map[string]interface{}` | Key-value pairs to include in Intent extras |

---

## Android Requirements

### Static Receivers (Manifest)

If you need to receive broadcasts while the app is closed (e.g., Boot Completion), add the receiver to your `AndroidManifest.xml`:

```xml
<uses-permission android:name="android.permission.RECEIVE_BOOT_COMPLETED" />

<application>
    <receiver android:name="com.sweetjuice.pkg.broadcast.BootReceiver" android:exported="false">
        <intent-filter>
            <action android:name="android.intent.action.BOOT_COMPLETED" />
        </intent-filter>
    </receiver>
</application>
```

---

## Notes

- Android only.
- Intent extras are converted to a generic Go `map[string]interface{}`.
- For sensitive system events, ensure your app has the required permissions in the manifest.
