import Foundation
import Sweetjuice

@_cdecl("NativeRenderUI")
public func NativeRenderUI(json: UnsafePointer<Int8>) {
    let jsonString = String(cString: json)
    UIManager.shared.updateUI(json: jsonString)
}

@_cdecl("NativeCallPlatform")
public func NativeCallPlatform(domain: UnsafePointer<Int8>, action: UnsafePointer<Int8>, args: UnsafePointer<Int8>) -> UnsafePointer<Int8>? {
    let d = String(cString: domain)
    let a = String(cString: action)
    let ar = String(cString: args)

    var response = "{}"

    switch d {
    case "system":
        response = SystemPlugin.shared.handleAction(a, args: ar)
    case "app":
        response = AppPlugin.shared.handleAction(a, args: ar)
    case "broadcast":
        response = BroadcastPlugin.shared.handleAction(a, args: ar)
    case "workmanager":
        response = WorkManagerPlugin.shared.handleAction(a, args: ar)
    case "notification":
        response = NotificationPlugin.shared.handleAction(a, args: ar)
    case "notification-listener":
        response = NotificationListenerPlugin.shared.handleAction(a, args: ar)
    default:
        print("NativeCallPlatform: Unknown domain \(d)")
    }

    return (response as NSString).utf8String
}
