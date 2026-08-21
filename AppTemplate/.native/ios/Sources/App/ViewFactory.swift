import SwiftUI

struct ViewFactory {
    @ViewBuilder
    static func build(node: JSONNode) -> some View {
        let view = internalBuild(node: node)
        applyStyles(view: AnyView(view), style: node.style)
    }

    @ViewBuilder
    private static func internalBuild(node: JSONNode) -> some View {
        switch node.type {
        case "column", "mu3:vstack":
            VStack {
                if let children = node.children {
                    ForEach(children, id: \.id) { child in
                        build(node: child)
                    }
                }
            }
        case "row", "mu3:hstack":
            HStack {
                if let children = node.children {
                    ForEach(children, id: \.id) { child in
                        build(node: child)
                    }
                }
            }
        case "mu3:text":
            Text(node.value ?? "")
        case "mu3:button":
            Button(action: {
                // Send click event to Go
                let event: [String: Any] = [
                    "id": node.id,
                    "name": "click",
                    "data": [:]
                ]
                if let data = try? JSONSerialization.data(withJSONObject: event),
                   let json = String(data: data, encoding: .utf8) {
                    JuiceappHandleNativeAction("ui:event", json)
                }
            }) {
                Text(node.text ?? "")
                    .frame(maxWidth: .infinity)
            }
            .buttonStyle(.borderedProminent)
        case "spacer":
            Spacer()
        case "mu3:box":
             ZStack {
                if let children = node.children {
                    ForEach(children, id: \.id) { child in
                        build(node: child)
                    }
                }
            }
            .padding()
            .background(Color.white)
            .cornerRadius(12)
            .shadow(radius: 4)
        default:
            Text("Unknown: \(node.type)")
        }
    }

    private static func applyStyles(view: AnyView, style: [String: AnyCodable]?) -> some View {
        var v = view

        if let style = style {
            if let colorHex = style["backgroundColor"]?.value as? String {
                v = AnyView(v.background(Color(hex: colorHex)))
            }
            if let radius = style["cornerRadius"]?.value as? Double {
                v = AnyView(v.cornerRadius(CGFloat(radius)))
            }
            if let padding = style["padding"]?.value as? Double {
                v = AnyView(v.padding(CGFloat(padding)))
            }
        }

        return v
    }
}

extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3: // RGB (12-bit)
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6: // RGB (24-bit)
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8: // ARGB (32-bit)
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (1, 1, 1, 0)
        }
        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue: Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}
