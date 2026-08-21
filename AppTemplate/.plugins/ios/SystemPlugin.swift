import UIKit

public class SystemPlugin {
    public static let shared = SystemPlugin()

    public func handleAction(_ action: String, args: String) -> String {
        if action == "getInfo" {
            let device = UIDevice.current
            let info: [String: Any] = [
                "model": device.model,
                "name": device.name,
                "systemName": device.systemName,
                "systemVersion": device.systemVersion,
                "platform": "ios"
            ]
            if let data = try? JSONSerialization.data(withJSONObject: info),
               let json = String(data: data, encoding: .utf8) {
                return json
            }
        }
        return "{}"
    }
}
