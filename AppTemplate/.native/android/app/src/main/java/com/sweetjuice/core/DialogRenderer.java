package com.sweetjuice.core;

import android.content.Context;
import android.util.Log;
import com.google.android.material.dialog.MaterialAlertDialogBuilder;
import org.json.JSONObject;
import juiceapp.Juiceapp;

/**
 * DialogRenderer facilitates the display of native Android Material 3 dialogs
 * triggered by the Go backend via the {@code ui:dialog} widget type.
 */
public class DialogRenderer {
    private static final String TAG = "SweetJuice";
    private final Context context;

    public DialogRenderer(Context context) {
        this.context = context;
    }

    /**
     * Renders and displays a native Android Material 3 Alert Dialog 
     * from a Sweet Juice dialog node.
     * 
     * @param node the JSON representation of the dialog.
     */
    public void showNativeDialog(JSONObject node) {
        try {
            String title = node.optString("title", "");
            String message = node.optString("message", "");
            String confirmText = node.optString("confirmText", "OK");
            String cancelText = node.optString("cancelText", "");
            String id = node.optString("id", "");

            MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(context)
                .setTitle(title)
                .setMessage(message)
                .setPositiveButton(confirmText, (dialog, which) -> {
                    sendEvent(id, "confirm", confirmText);
                });

            if (!cancelText.isEmpty()) {
                builder.setNegativeButton(cancelText, (dialog, which) -> {
                    sendEvent(id, "cancel", cancelText);
                });
            }

            builder.setCancelable(false)
                .show();
        } catch (Exception e) {
            Log.e(TAG, "UIManager: showNativeDialog failed", e);
        }
    }

    private void sendEvent(String id, String name, String buttonLabel) {
        try {
            JSONObject event = new JSONObject();
            event.put("id", id);
            event.put("name", name);
            JSONObject data = new JSONObject();
            data.put("button", buttonLabel);
            event.put("data", data);
            Juiceapp.handleMessageFromFrontend("ui:event", event.toString());
        } catch (Exception e) {
            Log.e(TAG, "UIManager: dialog event dispatch failed", e);
        }
    }
}
