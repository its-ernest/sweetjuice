# SweetJuice Project State Summary

## 1. Architectural Progress
- **Hybrid Bridge**: Successfully established a Go-centric Tri-Bridge architecture. Business logic and UI declarations live in Go, while rendering is delegated to Native (Android Views / SwiftUI).
- **Package Resiliency**: Core framework classes moved to `com.sweetjuice.core` (Android) and `Sources/App` (iOS) to support dynamic package renaming via `config.ini`.
- **Managed Directories**: 
    - `.native/`: Contains stable platform templates (not edited by user).
    - `.plugins/`: Source of truth for native plugin code, synchronized to the build folder by the CLI.
    - `app_assets/`: Static assets synchronized to native asset folders.

## 2. Recent Features & Fixes
- **Overlay System**: Flutter-inspired layered UI. Supports custom Go-driven boxes with background dimming on both Android and iOS.
- **Material 3 Dialogs**: High-level `ui.Dialog` API in Go that triggers native `MaterialAlertDialogBuilder` on Android.
- **Android Bridge Stabilization**: Re-exported all JNI-critical functions (`HandleMessageFromFrontend`, `SetNativeCallHandler`) in the `juiceapp` package to fix `gobind` visibility issues.
- **Smart CLI**: Android build now skips custom APK signing in `debug` mode for faster iteration.

## 3. Current Task: iOS Cross-Compilation
We have initiated the iOS transition. The SwiftUI renderer is in place, but since we are on Linux, we need a CI-based build pipeline.

### Next Steps (Unfinished):
1.  **`juice-cross` Integration**: Finalize the sibling repository logic.
2.  **`juice --run-cross ios`**:
    - Clone `juice-cross` into a temporary workspace.
    - Sync the current Go codebase and `.plugins/ios` into the workspace.
    - Push to the user's GitHub fork.
    - Implement a polling mechanism to wait for the GitHub Action to complete.
    - Download the built `Sweetjuice.xcframework.zip` (using `nightly.link` or GH API).
    - Extract and move to `.native/ios/` for final packaging via `xtool`.
3.  **iOS Plugin Parity**: Stabilize `Broadcast`, `WorkManager`, and `Notification` plugins for the iOS platform.

## 4. Technical Constants
- **Core Version**: v1.4.0
- **Android Target**: API 21+ (Min), 34 (Target)
- **iOS Target**: iOS 14.0+ (SwiftUI)
