package com.sweetjuice.app;

import android.content.Context;
import android.graphics.Color;
import android.graphics.Typeface;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.TextView;
import com.google.android.material.button.MaterialButton;
import com.google.android.material.card.MaterialCardView;
import org.json.JSONArray;
import org.json.JSONObject;
import java.util.HashSet;
import java.util.Set;
import sweetjuice.Sweetjuice;

public class UIManager {
    private final Context context;
    private final ViewGroup rootContainer;

    public UIManager(Context context, ViewGroup rootContainer) {
        this.context = context;
        this.rootContainer = rootContainer;
    }

    public void render(String jsonTree) {
        try {
            JSONObject rootNode = new JSONObject(jsonTree);
            View existingRoot = rootContainer.getChildAt(0);
            
            View newView = updateOrCreateView(existingRoot, rootNode);
            
            if (existingRoot != newView) {
                rootContainer.removeAllViews();
                rootContainer.addView(newView);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private View updateOrCreateView(View existingView, JSONObject node) throws Exception {
        String id = node.optString("id");
        String type = node.optString("type");
        
        View view = existingView;
        // If type or id mismatch, we must create a new view
        if (view == null || !id.equals(view.getTag()) || !isViewTypeOf(view, type)) {
            view = createView(type);
            view.setTag(id);
        }

        updateView(view, node);
        return view;
    }

    private boolean isViewTypeOf(View view, String type) {
        switch (type) {
            case "column": case "row": return view instanceof LinearLayout && !(view instanceof MaterialCardView);
            case "card": return view instanceof MaterialCardView;
            case "text": return view instanceof TextView && !(view instanceof MaterialButton) && !(view instanceof EditText);
            case "button": return view instanceof MaterialButton;
            case "textfield": return view instanceof EditText;
            case "spacer": return view.getClass().equals(View.class);
            default: return false;
        }
    }

    private View createView(String type) {
        switch (type) {
            case "column":
                LinearLayout col = new LinearLayout(context);
                col.setOrientation(LinearLayout.VERTICAL);
                return col;
            case "row":
                LinearLayout row = new LinearLayout(context);
                row.setOrientation(LinearLayout.HORIZONTAL);
                return row;
            case "card":
                MaterialCardView card = new MaterialCardView(context);
                LinearLayout cardLayout = new LinearLayout(context);
                cardLayout.setOrientation(LinearLayout.VERTICAL);
                card.addView(cardLayout);
                return card;
            case "text":
                return new TextView(context);
            case "button":
                return new MaterialButton(context);
            case "textfield":
                return new EditText(context);
            case "spacer":
                return new View(context);
            default:
                return new View(context);
        }
    }

    private void updateView(View view, JSONObject node) throws Exception {
        String type = node.optString("type");
        JSONObject style = node.optJSONObject("style");
        JSONArray events = node.optJSONArray("events");
        Set<String> eventSet = new HashSet<>();
        if (events != null) {
            for (int i = 0; i < events.length(); i++) eventSet.add(events.getString(i));
        }

        // Handle Children for Layouts
        if (view instanceof ViewGroup) {
            ViewGroup container = (view instanceof MaterialCardView) ? (ViewGroup)((MaterialCardView)view).getChildAt(0) : (ViewGroup)view;
            JSONArray children = node.optJSONArray("children");
            if (children != null) {
                // Basic reconciliation for children
                int nodeCount = children.length();
                int viewCount = container.getChildCount();
                
                // Update or add children
                for (int i = 0; i < nodeCount; i++) {
                    JSONObject childNode = children.getJSONObject(i);
                    View existingChild = (i < viewCount) ? container.getChildAt(i) : null;
                    View updatedChild = updateOrCreateView(existingChild, childNode);
                    
                    if (existingChild == null) {
                        container.addView(updatedChild);
                    } else if (existingChild != updatedChild) {
                        container.removeViewAt(i);
                        container.addView(updatedChild, i);
                    }
                }
                // Remove excess views
                if (viewCount > nodeCount) {
                    container.removeViews(nodeCount, viewCount - nodeCount);
                }
            } else {
                container.removeAllViews();
            }
        }

        // Specific View Properties
        if (view instanceof TextView && !(view instanceof EditText)) {
            String val = node.optString("value");
            if (!val.equals(((TextView) view).getText().toString())) {
                ((TextView) view).setText(val);
            }
        } else if (view instanceof MaterialButton) {
            ((MaterialButton) view).setText(node.optString("text"));
        } else if (view instanceof EditText) {
            EditText et = (EditText) view;
            String val = node.optString("value");
            // Only update text if it's different to preserve cursor position
            if (!val.equals(et.getText().toString())) {
                et.setText(val);
            }
            et.setHint(node.optString("placeholder"));
        } else if ("spacer".equals(type)) {
            int w = dpToPx((float) node.optDouble("width", 0));
            int h = dpToPx((float) node.optDouble("height", 0));
            view.setLayoutParams(new LinearLayout.LayoutParams(w == 0 ? 0 : w, h == 0 ? 0 : h, (w == 0 && h == 0) ? 1.0f : 0));
        }

        setupEvents(view, node.optString("id"), eventSet);
        applyStyles(view, type, style);
    }

    private void setupEvents(View view, final String id, Set<String> events) {
        // Remove existing listeners if possible or just overwrite
        if (events.contains("click")) {
            view.setOnClickListener(v -> sendEvent(id, "click", null));
        } else {
            view.setOnClickListener(null);
        }

        if (events.contains("long_press")) {
            view.setOnLongClickListener(v -> {
                sendEvent(id, "long_press", null);
                return true;
            });
        } else {
            view.setOnLongClickListener(null);
        }

        if (view instanceof EditText) {
            EditText et = (EditText) view;
            // Remove old watcher if any (tag-based storage)
            TextWatcher oldWatcher = (TextWatcher) et.getTag(R.id.text_watcher_tag);
            if (oldWatcher != null) et.removeTextChangedListener(oldWatcher);

            if (events.contains("changed")) {
                TextWatcher watcher = new TextWatcher() {
                    @Override public void beforeTextChanged(CharSequence s, int start, int count, int after) {}
                    @Override public void onTextChanged(CharSequence s, int start, int before, int count) {}
                    @Override public void afterTextChanged(Editable s) {
                        // Avoid loop: only send if the change is from user, not from updateView
                        if (et.hasFocus()) {
                            sendEvent(id, "changed", s.toString());
                        }
                    }
                };
                et.addTextChangedListener(watcher);
                et.setTag(R.id.text_watcher_tag, watcher);
            }
        }
    }

    private void applyStyles(View view, String type, JSONObject style) {
        if (style == null) return;

        if (style.has("flex")) {
            if (view.getLayoutParams() instanceof LinearLayout.LayoutParams) {
                LinearLayout.LayoutParams lp = (LinearLayout.LayoutParams) view.getLayoutParams();
                lp.weight = (float) style.optDouble("flex");
                if (lp.weight > 0) {
                    if (view.getParent() instanceof LinearLayout && ((LinearLayout)view.getParent()).getOrientation() == LinearLayout.HORIZONTAL) {
                        lp.width = 0;
                        lp.height = ViewGroup.LayoutParams.MATCH_PARENT;
                    } else {
                        lp.height = 0;
                        lp.width = ViewGroup.LayoutParams.MATCH_PARENT;
                    }
                }
            }
        }

        int p = dpToPx((float) style.optDouble("padding", 0));
        int pv = dpToPx((float) style.optDouble("paddingVertical", 0));
        int ph = dpToPx((float) style.optDouble("paddingHorizontal", 0));
        
        int pt = pv != 0 ? pv : p;
        int pb = pv != 0 ? pv : p;
        int pl = ph != 0 ? ph : p;
        int pr = ph != 0 ? ph : p;

        if (view instanceof MaterialCardView) {
            ((MaterialCardView)view).setContentPadding(pl, pt, pr, pb);
        } else {
            view.setPadding(pl, pt, pr, pb);
        }

        if (style.has("backgroundColor")) {
            try {
                int color = Color.parseColor(style.optString("backgroundColor"));
                if (view instanceof MaterialCardView) {
                    ((MaterialCardView) view).setCardBackgroundColor(color);
                } else {
                    view.setBackgroundColor(color);
                }
            } catch (Exception ignored) {}
        }

        if (style.has("cornerRadius")) {
            float radius = dpToPx((float) style.optDouble("cornerRadius"));
            if (view instanceof MaterialCardView) {
                ((MaterialCardView) view).setRadius(radius);
            }
        }

        if (view instanceof LinearLayout) {
            LinearLayout ll = (LinearLayout) view;
            int gravity = 0;
            if ("center".equals(style.optString("alignItems"))) gravity |= Gravity.CENTER_HORIZONTAL;
            if ("center".equals(style.optString("justifyContent"))) gravity |= Gravity.CENTER_VERTICAL;
            ll.setGravity(gravity);
        }

        if (view instanceof TextView) {
            TextView tv = (TextView) view;
            if (style.has("fontSize")) tv.setTextSize(TypedValue.COMPLEX_UNIT_SP, (float) style.optDouble("fontSize"));
            if (style.has("color")) {
                try { tv.setTextColor(Color.parseColor(style.optString("color"))); } catch (Exception ignored) {}
            }
            if ("bold".equals(style.optString("weight"))) tv.setTypeface(null, Typeface.BOLD);
        }
    }

    private int dpToPx(float dp) {
        return (int) TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, dp, context.getResources().getDisplayMetrics());
    }

    private void sendEvent(String id, String name, Object data) {
        try {
            JSONObject event = new JSONObject();
            event.put("id", id);
            event.put("name", name);
            if (data != null) event.put("data", data);
            Sweetjuice.handleMessageFromFrontend("ui:event", event.toString());
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
