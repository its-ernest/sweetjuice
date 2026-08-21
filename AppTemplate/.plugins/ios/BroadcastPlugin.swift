import Foundation

public class BroadcastPlugin {
    public static let shared = BroadcastPlugin()

    public func handleAction(_ action: String, args: String) -> String {
        switch action {
        case "register":
            return "{}"
        case "send":
            return "{}"
        case "post":
            return "{\"status\":\"ok\"}"
        default:
            return "{\"error\":\"Unknown broadcast action\"}"
        }
    }
}
