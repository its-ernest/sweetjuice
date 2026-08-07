package com.sweetjuice.app;

import android.app.AlertDialog;
import android.content.Context;
import android.util.Log;
import org.json.JSONObject;
import juiceapp.Juiceapp;

/**
 * Renders native {@link AlertDialog} overlays for {@code ui:dialog} nodes.
 *
 * <p>The dialog is non-cancelable and emits a {@code confirm} event back to Go
 * when the positive button is tapped.</p>
 */
class DialogRenderer {
    private static final String TAG = "SweetJuice";
    private final Context context;

    DialogRenderer(Context context) {
        this.context = context;
    }

    /**
     * Shows a native alert dialog from a {@code ui:dialog} JSON node.
     *
     * <p>Expected node fields:</p>
     * <ul>
     *   <li>{@code title} &rarr; dialog title</li>
     *   <li>{@code message} &rarr; dialog body text</li>
     *   <li>{@code buttonText} &rarr; positive button label, defaults to {@code OK}</li>
     *   <li>{@code id} &rarr; widget id echoed back in the confirm event</li>
     * </ul>
     *
     * <p>On confirm, the following payload is sent to Go:
     * {@code {"id":"<widgetId>","name":"confirm","data":{"button":"<buttonText>"}}}</p>
     *
     * @param node the dialog JSON node
     */
    void showNativeDialog(JSONObject node) {
        try {
            String title = node.optString("title", "");
            String message = node.optString("message", "");
            String buttonText = node.optString("buttonText", "OK");
            String id = node.optString("id", "");

            new AlertDialog.Builder(context)
                .setTitle(title)
                .setMessage(message)
                .setPositiveButton(buttonText, (dialog, which) -> {
                    try {
                        JSONObject event = new JSONObject();
                        event.put("id", id);
                        event.put("name", "confirm");
                        JSONObject data = new JSONObject();
                        data.put("button", buttonText);
                        event.put("data", data);
                        String payload = event.toString();
                        Juiceapp.handleMessageFromFrontend("ui:event", payload);
                    } catch (Exception e) {
                        Log.e(TAG, "UIManager: dialog confirm event failed", e);
                    }
                })
                .setCancelable(false)
                .show();
        } catch (Exception e) {
            Log.e(TAG, "UIManager: showNativeDialog failed", e);
        }
    }
}
