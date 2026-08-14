# App Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The App plugin handles application lifecycle events. It emits events when the app is resumed from the native side.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/app"
)

func init() {
    appPlugin := app.NewAppPlugin()
    appPlugin.Init(app)
}
```

---

## API Reference

### `app.NewAppPlugin() *AppPlugin`

Creates a new App plugin instance.

---

### `(*AppPlugin).Init(app *core.Application) error`

Initializes the plugin. Registers the `app:resumed` native callback.

---

### `(*AppPlugin).Name() string`

Returns the plugin name: `"app"`.

---

### `(*AppPlugin).Domain() string`

Returns the plugin domain: `"app"`.

---

## Events

### `app:resumed`

Fired when the app is resumed from the native side.

---

## Notes

- This plugin is typically included by default in AppTemplate projects.
- No special manifest permissions required.
