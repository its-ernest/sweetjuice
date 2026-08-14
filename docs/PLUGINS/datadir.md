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
    dirs, err := datadir.NewPlugin().GetDirs()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Files dir:", dirs.Files)
    fmt.Println("Cache dir:", dirs.Cache)
    fmt.Println("External files:", dirs.ExternalFiles)
}

func fileOps() {
    plugin := datadir.NewPlugin()
    
    // Write a file
    err := plugin.WriteFile("test.txt", "Hello SweetJuice")
    
    // Read a file
    content, err := plugin.ReadFile("test.txt")
    
    // Check existence
    exists, err := plugin.FileExists("test.txt")
    
    // Delete
    err = plugin.DeleteFile("test.txt")
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
    Files         string // Internal files directory
    Cache         string // Internal cache directory
    ExternalFiles string // External files directory (if available)
    ExternalCache string // External cache directory (if available)
}
```

---

### `(*DataDirPlugin).ReadFile(path string) (string, error)`

Reads a file from the app's internal files directory.

---

### `(*DataDirPlugin).WriteFile(path string, content string) error`

Writes content to a file in the app's internal files directory.

---

### `(*DataDirPlugin).FileExists(path string) (bool, error)`

Checks whether a file exists in the app's internal files directory.

---

### `(*DataDirPlugin).DeleteFile(path string) error`

Deletes a file from the app's internal files directory.

---

## Notes

- Paths are platform-specific and should not be hardcoded.
- External directories may not be available on all devices or if external storage is not mounted.
