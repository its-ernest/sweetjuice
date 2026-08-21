import Foundation

public class NotificationListenerPlugin {
    public static let shared = NotificationListenerPlugin()

    public func handleAction(_ action: String, args: String) -> String {
        return "{\"status\":\"processed\"}"
    }
}
