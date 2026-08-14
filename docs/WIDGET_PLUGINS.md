# SweetJuice Widget Plugin Design

## Core Concept

Widget plugins are third-party extensions that register **custom node types** with the native renderer. A plugin can:
- Declare a new widget type (e.g. `mu3:filled-card`)
- Provide a native view factory that the framework calls during reconciliation
- Bind events to Go handlers using the existing event dispatch system

## Plugin Contract

### Android Side

```java
public interface SweetJuiceWidgetFactory {
    String getType();
    View createView(Context ctx, JSONObject node, ViewGroup parent);
    void updateView(View view, JSONObject node);
}
```

Plugins register factories through the existing `SweetJuicePlugin` interface:

```java
public class Mu3Plugin implements SweetJuicePlugin {
    @Override
    public String getDomain() { return "mu3"; }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if ("widget:register".equals(action)) {
            // register factory with UIManager
        }
        return "{}";
    }
}
```

### Go Side

```go
type WidgetNode struct {
    BaseNode
    Props map[string]interface{}
}

func (n *WidgetNode) Serialize() (map[string]interface{}, error) {
    return map[string]interface{}{
        "type": n.Type,          // e.g. "mu3:card"
        "id":   n.BaseNode.ID,
        "props": n.Props,
        "events": n.Events,
    }, nil
}
```

Usage:

```go
ui.NewWidget("mu3:card").
    ID("my_card").
    Prop("title", "Hello").
    Prop("elevation", 2.0).
    OnClick(func() { ... })
```

## UIManager Extension

`UIManager` delegates unknown types to a `widgetFactories` map:

```java
private final Map<String, SweetJuiceWidgetFactory> widgetFactories = new HashMap<>();

public void registerWidgetFactory(SweetJuiceWidgetFactory factory) {
    widgetFactories.put(factory.getType(), factory);
}
```

During `updateOrCreateView`, if type doesn’t match builtins:

```java
SweetJuiceWidgetFactory factory = widgetFactories.get(type);
if (factory != null) {
    return factory.createView(context, node, container);
}
```

## Event Integration

Widgets emit events through the existing `sendEvent(id, name, data)` path. No changes needed.

## Kotlin Enablement

Add to `build.gradle`:

```gradle
plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
}

kotlin { jvmToolchain(17) }

dependencies {
    implementation libs.androidx.core.ktx
    implementation libs.androidx.lifecycle.runtime.ktx
    // Material 3 Compose or Material3-lite
    implementation "com.google.android.material:material:<MATERIAL3_VERSION>"
}
```

## Future: Compose Widgets

For full Compose rendering, the factory contract can extend:

```java
@Composable
fun RenderWidget(type: String, props: Map<String, Any>, onEvent: (String, Any?) -> Unit)
```

This lets widget plugins ship Compose-based renderables without blocking the existing View reconciler.
