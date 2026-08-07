package com.sweetjuice.app;

import com.sweetjuice.plugin.SweetJuiceWidgetFactory;
import java.util.HashMap;
import java.util.Map;

/**
 * Maintains the registry of widget factories keyed by type string.
 *
 * <p>UIManager queries this registry during view reconciliation to decide
 * whether a node should be delegated to a plugin-provided factory.</p>
 */
class WidgetRegistry {
    private final Map<String, SweetJuiceWidgetFactory> widgetFactories = new HashMap<>();

    /**
     * Registers a widget factory using the type returned by {@link SweetJuiceWidgetFactory#getType()}.
     *
     * @param factory the factory to register; its type string is used as the lookup key
     */
    void registerWidgetFactory(SweetJuiceWidgetFactory factory) {
        widgetFactories.put(factory.getType(), factory);
    }

    /**
     * Registers a widget factory under an explicit type string.
     *
     * @param type    the widget type identifier to register under
     * @param factory the factory to register
     */
    void registerWidgetFactory(String type, SweetJuiceWidgetFactory factory) {
        widgetFactories.put(type, factory);
    }

    /**
     * Looks up a registered widget factory by type.
     *
     * @param type the widget type identifier
     * @return the matching factory, or {@code null} if no factory is registered for the type
     */
    SweetJuiceWidgetFactory get(String type) {
        return widgetFactories.get(type);
    }
}
