# Mu3 Plugin

> **Status:** Stable  
> **Platforms:** Android  
> **Version:** v1 alpha

The Mu3 plugin provides Material You (Material 3) UI components for building native Android interfaces from Go.

---

## Quick Example

```go
package main

import (
    "github.com/sweet-juice/sweetjuice/plugins/mu3"
    "github.com/sweet-juice/sweetjuice/ui"
)

func renderHome() ui.Node {
    return ui.Root(
        ui.VStack(
            mu3.TopAppBar("Home", ""),
            mu3.Button("Tap Me").OnClick(func() {
                fmt.Println("Clicked")
            }),
            mu3.Card("Title", "Subtitle"),
        ),
        "#FFFBFE",
    )
}
```

---

## Layout

| Function | Description |
|----------|-------------|
| `VStack(children ...ui.Node)` | Vertical layout |
| `HStack(children ...ui.Node)` | Horizontal layout |
| `Spacer()` | Flexible empty space |
| `Root(child, backgroundColor)` | Root container with background color |

---

## Buttons

| Function | Description |
|----------|-------------|
| `Button(text)` | Filled button |
| `TextButton(text)` | Text-only button |
| `OutlinedButton(text)` | Outlined button |
| `TonalButton(text)` | Tonal button |
| `ElevatedButton(text)` | Elevated button |
| `IconButton(name)` | Icon button |
| `SegmentedButton(options, selected)` | Segmented control |
| `ButtonGroup(children)` | Button group |

---

## Input & Media

| Function | Description |
|----------|-------------|
| `TextField(placeholder)` | Text input field |
| `Image(src)` | Image from URI |
| `Video(src)` | Video player |

---

## Containers

| Function | Description |
|----------|-------------|
| `Card(title, subtitle)` | Material card |
| `Dialog(title, message, buttonText)` | Alert dialog |

---

## Navigation

| Function | Description |
|----------|-------------|
| `TopAppBar(title, navigationIcon)` | Top app bar |
| `BottomAppBar()` | Bottom app bar |
| `NavigationBar(destinations)` | Bottom navigation bar |
| `NavigationRail(destinations)` | Side navigation rail |
| `SearchBar(hint, showDocked)` | Search bar |
| `Tabs(primary, secondary, selected)` | Tab layout |
| `Toolbar(items, orientation)` | Toolbar |

---

## Events

All widgets support:
- `.OnClick(func())`
- `.OnLongPress(func())`
- `.OnChanged(func(string))` for input widgets

---

## Notes

- Android only. Renders native Material 3 components.
- Icons use Material Symbols font.
