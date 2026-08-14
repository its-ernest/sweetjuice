# DataStore Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The DataStore plugin provides a simple key-value string store backed by the native platform's persistent storage.

---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/sweet-juice/sweetjuice/plugins/datastore"
)

func savePrefs() {
    datastore.Set("username", "alice")
    datastore.Set("theme", "dark")
}

func loadPrefs() {
    username, _ := datastore.Get("username", "guest")
    theme, _ := datastore.Get("theme", "light")
    fmt.Println("User:", username, "Theme:", theme)
}

func clearPrefs() {
    datastore.Clear()
}
```

---

## API Reference

### `datastore.NewPlugin() *DataStorePlugin`

Creates a new datastore plugin instance.

---

### `(*DataStorePlugin).Set(key, value string) error`

Stores a string value.

---

### `(*DataStorePlugin).Get(key, defaultVal string) (string, error)`

Retrieves a string value. Returns `defaultVal` if the key does not exist.

---

### `(*DataStorePlugin).Delete(key string) error`

Deletes a key-value pair.

---

### `(*DataStorePlugin).Clear() error`

Removes all stored data.

---

### `(*DataStorePlugin).GetAll() (map[string]string, error)`

Returns all key-value pairs.

---

## Notes

- Values are stored as strings. For complex data, serialize to JSON.
- Data persists across app restarts.
- No special permissions required.
