# Sweet Juice — Compact Context

## Project
Go-based mobile app framework (`github.com/sweet-juice/sweetjuice`). Builds native Android apps from Go via declarative UI + gomobile bindings. CLI: `juice`.

## Architecture
- **Go Engine**: State, event bus, declarative UI serialization
- **Native Shell**: Android (Java/Gradle) renders JSON→views, handles system APIs
- **Bridge**: JSON strings with `"domain:action"` routing

## Key Packages
- `core/`: Runtime (`app.go`, `bridge.go`, `eventbus.go`)
- `app/`: Bootstrap (`Run`, `ReRender`)
- `ui/`: Declarative UI nodes + Material 3 widgets
- `plugins/`: Built-ins (permission, notification, biometric, datastore, etc.)
- `AppTemplate/`: Project scaffold copied by `juice --new`

## AppTemplate Structure
- `engine.go`: Exports `StartApplication()`, `HandleMessageFromFrontend()`, etc. Package: `juiceapp`
- `plugins.go`: `registerPluginDefinitions()` (before Run) + `initPlugins()` (after Run)
- `start.go`: Bootstrap, creates state + root view
- `lib/state/`: Example state with waiting/dialog/tab management
- `lib/views/`: Example view with permission buttons
- `native/android/`: Full Android Studio project

## Communication Flow
```
Go → Native: CallNativePlatform("domain:action", json)
Native → Go: handleNativeAction("domain:callback", json)
UI Events: handleMessageFromFrontend("ui:event", json)
Rendering: app.ReRender() → Serialize() → JSON → ui:render
```

## Plugin Registration Flow
1. `registerPluginDefinitions()` — calls `RegisterPlugin()` for all plugins BEFORE `app.Run()`
2. `app.Run(root)` — initializes global app, triggers first render
3. Java reads `GetRegisteredPlugins()`, registers widget factories
4. `initPlugins()` — calls `plugin.Init(app)` AFTER `app.Run()`
5. Java calls `Juiceapp.reRender()` for final render with all factories

## Current State
- Early stage, transitioning from webview to native components
- Android focus; iOS experimental
- Version v1.4.0 CLI, v1.0.0-alpha framework
- Active: `feature/sweetjuice-rebrand-and-plugins`

## Notable TODOs
- iOS Docker-Xtool integration
- Background services in Go
- Geolocation, Camera (separate repo)
- `deviceinfo`, `badge`, `share` plugins

## Recent Changes
- Dynamic plugin registration from Go side only
- Android derives widget factories from `GetRegisteredPlugins()`
- Permission plugin fixed: callback now handles JSON array payload
- Background location permission has its own button with fallback to app settings
- Special permissions plugin cleaned up (removed notification access)
