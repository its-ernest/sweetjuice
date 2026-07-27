# SweetJuice Project State Summary - July 2026

## Objective
Transitioning SweetJuice from a WebView-based wrapper to a 100% Go-only Native UI framework for Android (and eventually iOS).

## Current Architecture
- **Go Backend**: Manages state and declarative UI tree.
- **Native Layer (Android)**: Uses a `UIManager` to parse JSON UI trees and reconcile them into native Material Design components.
- **JNI Bridge**: Facilitates communication between Go and Java via `gobind`.

## Key Components
### 1. Core Engine (`core/`)
- `Application`: Global lifecycle and binding manager.
- `Bridge`: Handles `ui:event` (Native -> Go) and `ui:render` (Go -> Native).
- `SetUIEventDispatcher`: Routes native events to the UI registry.

### 2. UI Framework (`ui/`)
- `Node` interface for all UI elements.
- `BaseNode`: Handles common properties like stable `ID`, `Style`, and `Events`.
- **Widgets**: `Text`, `Button`, `TextField`, `VStack`, `HStack`, `Card`, `Spacer`.
- **Styling**: Modular `ui/style` package with Flexbox-like properties (`flex`, `padding`, `alignItems`, etc.).
- **Events**: Fluent API for `.OnClick()`, `.OnLongPress()`, and `.OnChanged()`.

### 3. Native Layout Engine (`UIManager.java`)
- **Stable Reconciliation**: Uses IDs to update existing views instead of rebuilding. Preserves focus and cursor position in text fields.
- **DP/SP Conversion**: Correctly handles device-independent units.
- **Circular Event Protection**: Only sends `changed` events if the view has focus.

### 4. Simplified Stateful API (`app/`)
- `Run(Component)`: entry point.
- `ReRender()`: Manual trigger for UI diff/update.
- `AppTemplate`: Refactored into a modular `lib/` structure:
    - `lib/state/`: Single source of truth with sub-states (e.g., `UserState`).
    - `lib/components/`: Reusable, state-aware UI units.
    - `lib/views/`: Coordinator components for layout.

## Recent Fixes
- **Blank Screen**: Resolved by adding explicit `reRender()` calls in `SweetJuiceActivity` and fixing `LayoutParams` in the native layer.
- **Lost Input**: Fixed by implementing stable IDs and a proper reconciliation algorithm in `UIManager`.
- **CLI Automation**: `juice --new` now leaves the `go.mod` as `module myapp` and internal imports as `myapp/` to match the template, avoiding broken paths.

## Next Steps
- Continue implementing more native widgets (Images, Lists, etc.).
- Expand styling capabilities.
- Implement iOS `UIManager` counterpart.
