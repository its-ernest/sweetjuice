package com.sweetjuice.app;

import android.text.Editable;
import android.text.TextWatcher;
import android.util.Log;
import android.view.View;
import org.json.JSONObject;
import java.util.Set;
import juiceapp.Juiceapp;

/**
 * Binds UI events from the widget tree onto native Android view listeners.
 *
 * <p>Supported event names:</p>
 * <ul>
 *   <li>{@code click} &rarr; {@link View.OnClickListener}</li>
 *   <li>{@code long_press} &rarr; {@link View.OnLongClickListener}</li>
 *   <li>{@code changed} &rarr; {@link TextWatcher} on {@link android.widget.EditText} only</li>
 * </ul>
 *
 * <p>Events are sent to the Go frontend as JSON payloads via
 * {@link Sweetjuice#handleMessageFromFrontend(String, String)} with the shape:
 * {@code {"id":"<widgetId>","name":"<eventName>","data":<optional>}}</p>
 */
class EventBinder {
    private static final String TAG = "SweetJuice";

    /**
     * Attaches or detaches event listeners on the view according to the requested event set.
     *
     * @param view   the native view to bind events onto
     * @param id     the widget id used in outgoing event payloads
     * @param events the set of event names requested by the node
     */
    void setupEvents(View view, String id, Set<String> events) {
        try {
            if (events.contains("click")) {
                Log.d(TAG, "UIManager: setting click listener for id=" + id);
                view.setOnClickListener(v -> {
                    Log.d(TAG, "UIManager: click fired for id=" + id);
                    sendEvent(id, "click", null);
                });
            } else {
                view.setOnClickListener(null);
            }
            if (events.contains("long_press")) {
                view.setOnLongClickListener(v -> {
                    Log.d(TAG, "UIManager: long_press fired for id=" + id);
                    sendEvent(id, "long_press", null);
                    return true;
                });
            } else {
                view.setOnLongClickListener(null);
            }
            if (view instanceof android.widget.EditText) {
                android.widget.EditText et = (android.widget.EditText) view;
                TextWatcher oldWatcher = (TextWatcher) et.getTag(R.id.text_watcher_tag);
                if (oldWatcher != null) et.removeTextChangedListener(oldWatcher);
                if (events.contains("changed")) {
                    TextWatcher watcher = new TextWatcher() {
                        @Override public void beforeTextChanged(CharSequence s, int start, int count, int after) {}
                        @Override public void onTextChanged(CharSequence s, int start, int before, int count) {}
                        @Override public void afterTextChanged(Editable s) {
                            if (et.hasFocus()) {
                                sendEvent(id, "changed", s.toString());
                            }
                        }
                    };
                    et.addTextChangedListener(watcher);
                    et.setTag(R.id.text_watcher_tag, watcher);
                }
            }
        } catch (Exception e) {
            Log.e(TAG, "UIManager: setupEvents failed", e);
        }
    }

    /**
     * Sends a UI event from native back to the Go frontend.
     *
     * @param id   the originating widget id
     * @param name the event name, e.g. {@code click}
     * @param data optional event data, may be {@code null}
     */
    void sendEvent(String id, String name, Object data) {
        try {
            JSONObject event = new JSONObject();
            event.put("id", id);
            event.put("name", name);
            if (data != null) event.put("data", data);
            String payload = event.toString();
            Log.d(TAG, "EventBinder: sendEvent id=" + id + " name=" + name + " payload=" + payload);
            Juiceapp.handleMessageFromFrontend("ui:event", payload);
            Log.d(TAG, "EventBinder: sendEvent completed for id=" + id);
        } catch (Exception e) {
            Log.e(TAG, "EventBinder: sendEvent failed", e);
            e.printStackTrace();
        }
    }
}
