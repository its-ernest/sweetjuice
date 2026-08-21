import UIKit
import Sweetjuice

class AppDelegate: NSObject, UIApplicationDelegate {
    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey : Any]? = nil) -> Bool {
        print("SweetJuice: iOS AppDelegate started")
        return true
    }

    // Standard lifecycle hooks for SweetJuice plugins to hook into
    func application(_ app: UIApplication, open url: URL, options: [UIApplication.OpenURLOptionsKey : Any] = [:]) -> Bool {
        // Deep linking support
        return true
    }
}
