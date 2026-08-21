package com.sweetjuice.pkg.mu3;

import android.content.Context;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.widget.LinearLayout;
import android.widget.TextView;
import androidx.appcompat.app.AppCompatActivity;
import com.google.android.material.dialog.MaterialAlertDialogBuilder;
import com.sweetjuice.core.SweetJuiceApp;
import com.sweetjuice.plugin.SweetJuiceWidgetFactory;
import org.json.JSONObject;
import juiceapp.Juiceapp;

public class Mu3DialogWidgetFactory implements SweetJuiceWidgetFactory {

    @Override
    public String getType() {
        return "mu3:dialog";
    }

    @Override
    public View createView(Context ctx, JSONObject node, ViewGroup parent) {
        AppCompatActivity activity = getActiveActivity(ctx);
        if (activity == null) {
            Log.w("Mu3Dialog", "No active activity for dialog");
            return new View(ctx);
        }

        String title = node.optString("title", "");
        String message = node.optString("message", "");
        String buttonText = node.optString("buttonText", "OK");
        String id = node.optString("id", "");

        new MaterialAlertDialogBuilder(activity)
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
                    Juiceapp.handleMessageFromFrontend("ui:event", event.toString());
                } catch (Exception e) {
                    Log.e("Mu3Dialog", "Failed to send confirm event", e);
                }
            })
            .setCancelable(false)
            .show();

        // Dialogs don't occupy space in the view hierarchy
        return new View(ctx);
    }

    @Override
    public void updateView(View view, JSONObject node) {
        // Dialogs are transient and shown immediately in createView
    }

    private AppCompatActivity getActiveActivity(Context ctx) {
        if (ctx.getApplicationContext() instanceof SweetJuiceApp) {
            return ((SweetJuiceApp) ctx.getApplicationContext()).getActiveActivity();
        }
        return null;
    }
}
