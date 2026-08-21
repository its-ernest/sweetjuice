import Foundation

public class WorkManagerPlugin {
    public static let shared = WorkManagerPlugin()

    public func handleAction(_ action: String, args: String) -> String {
        switch action {
        case "execute":
            return "{\"status\":\"success\"}"
        case "enqueueOneTime", "enqueueOneTimeWithDelay", "enqueuePeriodic":
            return "{\"error\":\"WorkManager not supported on iOS\"}"
        case "isEnqueued":
            return "{\"enqueued\":false,\"error\":\"not supported on ios\"}"
        case "cancelAll":
            return "{}"
        default:
            return "{\"error\":\"Unknown workmanager action\"}"
        }
    }
}
