# DataDir Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The DataDir plugin provides access to standard application directories, including files, cache, external files, and external cache paths.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/datadir"
)

func logDirs() {
    dirs, err := datadir.GetDirs()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Files dir:", dirs.FilesDir)
    fmt.Println("Cache dir:", dirs.CacheDir)
    fmt.Println("External files:", dirs.ExternalFilesDir)
}
```

---

## API Reference

### `datadir.NewPlugin() *DataDirPlugin`

Creates a new datadir plugin instance.

---

### `(*DataDirPlugin).GetDirs() (AppDirs, error)`

Returns the application directory paths.

**Returns:**

```go
type AppDirs struct {
    FilesDir        string // Internal files directory
    CacheDir        string // Internal cache directory
    ExternalFilesDir string // External files directory (if available)
    ExternalCacheDir string // External cache directory (if available)
}
```

---

## Notes

- Paths are platform-specific and should not be hardcoded.
- External directories may not be available on all devices or if external storage is not mounted.
