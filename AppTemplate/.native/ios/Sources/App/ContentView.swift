import SwiftUI
import Sweetjuice

struct ContentView: View {
    @StateObject private var uiManager = UIManager.shared

    var body: some View {
        ZStack {
            // Main Content
            VStack {
                if let root = uiManager.rootNode {
                    ViewFactory.build(node: root)
                } else {
                    ProgressView("Loading...")
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)

            // Overlay Layer
            if let overlay = uiManager.overlayNode {
                ZStack {
                    if overlay.dim {
                        Color.black.opacity(0.4)
                            .edgesIgnoringSafeArea(.all)
                    }
                    ViewFactory.build(node: overlay.child)
                }
            }
        }
        .onAppear {
            JuiceappReRender()
        }
    }
}
