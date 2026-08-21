import SwiftUI
import Sweetjuice

class UIManager: ObservableObject {
    static let shared = UIManager()

    @Published var rootNode: JSONNode?
    @Published var overlayNode: OverlayNode?

    func updateUI(json: String) {
        guard let data = json.data(using: .utf8) else { return }
        do {
            let node = try JSONDecoder().decode(JSONNode.self, from: data)
            DispatchQueue.main.async {
                if node.type == "ui:overlay" {
                    self.overlayNode = try? JSONDecoder().decode(OverlayNode.self, from: data)
                } else if node.type == "ui:overlay:dismiss" {
                    self.overlayNode = nil
                } else {
                    self.rootNode = node
                }
            }
        } catch {
            print("UIManager: Decode failed \(error)")
        }
    }
}

class JSONNode: Codable {
    let type: String
    let id: String
    let props: [String: AnyCodable]?
    let style: [String: AnyCodable]?
    let children: [JSONNode]?
    let value: String?
    let text: String?
    let child: JSONNode?
}

class OverlayNode: Codable {
    let type: String
    let id: String
    let child: JSONNode
    let dim: Bool
}

// Helper for AnyCodable to handle dynamic JSON from Go
struct AnyCodable: Codable {
    let value: Any

    init(_ value: Any) {
        self.value = value
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        if let x = try? container.decode(String.self) { value = x }
        else if let x = try? container.decode(Bool.self) { value = x }
        else if let x = try? container.decode(Int.self) { value = x }
        else if let x = try? container.decode(Double.self) { value = x }
        else { value = "" }
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        if let x = value as? String { try container.encode(x) }
        else if let x = value as? Bool { try container.encode(x) }
        else if let x = value as? Int { try container.encode(x) }
        else if let x = value as? Double { try container.encode(x) }
    }
}
