# FilePicker Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The FilePicker plugin triggers the native file picker UI for selecting files. Results are returned asynchronously via the `filepicker:result` event.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/filepicker"
)

func pickFile() {
    _, err := filepicker.PickFile(filepicker.PickerOptions{
        AllowMultiple: false,
        MimeTypes:     []string{"image/*", "application/pdf"},
    })
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

---

## API Reference

### `filepicker.NewPlugin() *FilePickerPlugin`

Creates a new file picker plugin instance.

---

### `(*FilePickerPlugin).PickFile(options PickerOptions) (string, error)`

Opens the native file picker.

| Field | Type | Description |
|-------|------|-------------|
| `AllowMultiple` | `bool` | Allow selecting multiple files |
| `MimeTypes` | `[]string` | Filter by MIME types, e.g. `["image/*"]` |

**Returns:** JSON string of selected file URIs, or error

---

## Events

### `filepicker:result`

Fired when the user selects a file or cancels.

Success:

```json
{
  "uri": "content://com.android.providers.media.documents/document/image%3A123",
  "name": "photo.jpg",
  "size": 102400
}
```

Cancelled:

```json
{
  "error": "cancelled"
}
```

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.READ_EXTERNAL_STORAGE" />
```

On Android 13+, scoped storage is used automatically.

---

## Notes

- Results are always delivered via event, not as a direct return value.
- On iOS, the picker supports cloud documents if configured.
