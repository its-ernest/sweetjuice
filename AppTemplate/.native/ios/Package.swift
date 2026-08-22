// swift-tools-version:5.5
import PackageDescription

let package = Package(
    name: "GenericApp",
    platforms: [
        .iOS(.v15)
    ],
    products: [
        .library(name: "App", targets: ["App"])
    ],
    dependencies: [],
    targets: [
        .target(
            name: "App",
            dependencies: [
                .target(name: "Sweetjuice")
            ],
            path: "Sources"
        ),
        .binaryTarget(
            name: "Sweetjuice",
            path: "Sweetjuice.xcframework"
        )
    ]
)
