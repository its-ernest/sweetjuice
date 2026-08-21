import SwiftUI
import Sweetjuice

@main
struct GenericAppApp: App {
    @UIApplicationDelegateAdaptor(AppDelegate.self) var delegate

    init() {
        // Initialize Go backend
        JuiceappStartApplication()
    }

    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
