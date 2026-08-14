# Deep Linking Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The Deep Linking plugin handles incoming deep links and dispatches them to registered handler functions in Go.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/deeplinking"
)

func init() {
    dl := deeplinking.NewPlugin()
    dl.OnURL(func(url string) {
        fmt.Println("Deep link received:", url)
        // Handle navigation based on URL
    })
}
```

---

## API Reference

### `deeplinking.NewPlugin() *DeepLinkingPlugin`

Creates a new deep linking plugin instance.

---

### `(*DeepLinkingPlugin).OnURL(handler URLHandler)`

Registers a handler function for incoming deep links.

```go
type URLHandler func(url string)
```

---

## Android Setup

Add intent filters to `AndroidManifest.xml`:

```xml
<activity android:name=".SweetJuiceActivity">
    <intent-filter>
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.DEFAULT" />
        <category android:name="android.intent.category.BROWSABLE" />
        <data android:scheme="https" android:host="example.com" />
    </intent-filter>
</activity>
```

---

## iOS Setup

Add URL schemes to `Info.plist`:

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>example</string>
        </array>
    </dict>
</array>
```

---

## Notes

- Deep links are delivered asynchronously via the registered handler.
- On Android, links must match the intent filter pattern.
- On iOS, the app must be configured with the correct URL scheme or universal link.
