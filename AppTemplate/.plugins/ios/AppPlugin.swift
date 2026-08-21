import Foundation
import Sweetjuice

public class AppPlugin {
    public static let shared = AppPlugin()

    public func handleAction(_ action: String, args: String) -> String {
        return "{}"
    }
}
