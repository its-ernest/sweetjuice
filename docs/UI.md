# Sweet Juice UI Framework

The UI framework allows you to build native Android and iOS interfaces using declarative Go code.

## Component Architecture

Everything in Sweet Juice UI is a **Component**. A component must implement the `Render()` method, which returns a tree of **Nodes**.

```go
type HomeView struct {}

func (v *HomeView) Render() ui.Node {
    return ui.VStack(
        ui.Text("Welcome Home"),
        ui.Button("Click Me"),
    )
}
```

## Standard Nodes

| Node | Purpose |
|---|---|
| `VStack` | Vertical layout container |
| `HStack` | Horizontal layout container |
| `Text` | Label for displaying text |
| `Button` | Standard clickable button |
| `TextField` | Input field for user text |
| `Spacer` | Flexible or fixed spacing |
| `Image` | Displays local or remote images |
| `Video` | Renders native video players |
| `Dialog` | Native Material 3 Alert Dialogs |

## Styling

Styles are defined using Go structs and converted to native attributes (DP/SP) automatically.

```go
ui.Text("Styled Text").Style(style.Text{
    FontSize: 16,
    Weight:   style.Bold,
    Color:    "#FF0000",
})
```

## Event Handling

Interactions are handled by passing Go closures to event methods.

```go
ui.Button("Submit").OnClick(func() {
    fmt.Println("Form submitted")
})
```

## Material 3 (mu3)

The `mu3` plugin provides advanced Material 3 widgets:
- `Box`: A flexible container with elevation and radius.
- `Card`: A title/subtitle based container.
- `TopAppBar`: Standard Android top bars.
- `NavigationBar`: Bottom navigation components.
