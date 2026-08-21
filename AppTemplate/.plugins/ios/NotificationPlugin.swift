import Foundation
import UserNotifications

public class NotificationPlugin {
    public static let shared = NotificationPlugin()
    private let center = UNUserNotificationCenter.current()

    public func handleAction(_ action: String, args: String) -> String {
        guard let data = args.data(using: .utf8),
              let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            return "{\"error\":\"Invalid JSON\"}"
        }

        switch action {
        case "post":
            center.requestAuthorization(options: [.alert, .sound]) { granted, error in
                if granted {
                    self.postNotification(json: json)
                }
            }
            return "{\"status\":\"posted\"}"
        case "cancel":
            if let id = json["id"] as? Int {
                center.removeDeliveredNotifications(withIdentifiers: ["\(id)"])
                return "{\"status\":\"cancelled\"}"
            }
            return "{\"error\":\"Missing id\"}"
        default:
            return "{\"error\":\"Unknown notification action\"}"
        }
    }

    private func postNotification(json: [String: Any]) {
        let id = json["id"] as? Int ?? 0
        let title = json["title"] as? String ?? ""
        let body = json["body"] as? String ?? ""
        let _ = json["channel_id"] as? String ?? "default_channel"
        let _ = json["channel_name"] as? String ?? "General Notifications"

        let content = UNMutableNotificationContent()
        content.title = title
        content.body = body
        content.sound = .default

        let trigger = UNTimeIntervalNotificationTrigger(timeInterval: 1, repeats: false)
        let request = UNNotificationRequest(identifier: "\(id)", content: content, trigger: trigger)
        center.add(request)
    }
}
