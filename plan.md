# Sweet Juice Development Plan

## Vision
A high-performance, Go-based mobile application framework using native components (via JSON-bridge) for maximum smoothness and low overhead, inspired by Flutter's declarative UI but leveraging React Native's native component approach.

## Core Objectives
1. **Framework Stability**: Ensure the Go engine, the JSON bridge, and the Native shells (Android/iOS) are robust and performant.
2. **Plugin Reliability**: 100% functional parity for all pre-built plugins across both platforms.
3. **Developer Experience (DX)**: Provide a seamless `AppTemplate` that stays in sync with framework evolution.
4. **Documentation Excellence**: Prioritize clarity over quantity. Every core component and plugin must be self-explanatory and well-documented.

## Roadmap

### Phase 1: Foundation & Plugin Audit
- [ ] **Plugin Validation**: Rigorous testing of all existing plugins (`biometric`, `datastore`, `gps`, `permissions`, `notifications`, `calls`, `sms`, etc.) on both Android and iOS.
- [ ] **Bridge Optimization**: Profile the Go <-> Native JSON communication to ensure zero lag.
- [ ] **Documentation Foundation**: Establish the "Clarity > Quantity" standard for all Go and Native source code.

### Phase 2: UI & Component Expansion
- [ ] **Mu3 (Material 3) Parity**: Ensure all Material 3 widgets in Go have perfect rendering and interaction in the Native shells.
- [ ] **New Plugins Implementation**:
    - `deviceinfo`: Hardware and software metadata.
    - `badge`: App icon notification badges.
    - `share`: Native sharing dialogs.
- [ ] **Refine UI Serialization**: Improve how complex UI trees are sent over the bridge.

### Phase 3: Ecosystem & Tooling
- [ ] **AppTemplate Synchronization**: Automate updates to the `AppTemplate` to prevent drift from the core framework.
- [ ] **iOS Deep Integration**: Complete the Docker-Xtool integration for iOS build pipelines.
- [ ] **Advanced Features**: Implement background services in Go and enhanced geolocation capabilities.

## Principles
- **Go-First**: All business logic, state management, and event handling must reside in Go.
- **Native Performance**: UI must use native widgets via the `WidgetRegistry`.
- **Declarative UI**: UI state -> JSON -> Native View.
- **Clarity over Quantity**: Code must be readable and follow idiomatic Go/Java/Swift patterns.
