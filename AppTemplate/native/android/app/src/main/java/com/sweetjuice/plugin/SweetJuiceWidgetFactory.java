package com.sweetjuice.plugin;

import android.content.Context;
import android.view.View;
import android.view.ViewGroup;
import org.json.JSONObject;

/**
 * SweetJuiceWidgetFactory is the extension point for third-party widget plugins.
 *
 * A plugin registers a factory for a custom widget type (e.g. "mu3:card").
 * During reconciliation, UIManager delegates unknown types to the matching factory.
 */
public interface SweetJuiceWidgetFactory {

    /**
     * Returns the widget type identifier (e.g. "mu3:card").
     */
    String getType();

    /**
     * Create a new native view for the given widget node.
     */
    View createView(Context ctx, JSONObject node, ViewGroup parent);

    /**
     * Update an existing view to match the new node properties.
     */
    void updateView(View view, JSONObject node);
}
