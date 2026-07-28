package com.sweetjuice.app;

import android.content.Context;
import android.content.res.ColorStateList;
import android.graphics.Color;
import android.graphics.Typeface;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.Log;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.TextView;
import androidx.core.content.ContextCompat;
import com.google.android.material.button.MaterialButton;
import com.google.android.material.button.MaterialButtonToggleGroup;
import com.google.android.material.card.MaterialCardView;
import com.google.android.material.floatingactionbutton.ExtendedFloatingActionButton;
import com.google.android.material.floatingactionbutton.FloatingActionButton;
import com.sweetjuice.plugin.SweetJuiceWidgetFactory;
import org.json.JSONArray;
import org.json.JSONObject;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import sweetjuice.Sweetjuice;

public class UIManager {
    private static final String TAG = "SweetJuice";
    private final Context context;
    private final ViewGroup rootContainer;
    private boolean renderFailed;
    private final Map<String, SweetJuiceWidgetFactory> widgetFactories = new HashMap<>();

    public UIManager(Context context, ViewGroup rootContainer) {
        this.context = context;
        this.rootContainer = rootContainer;
        this.renderFailed = false;
    }

    public void registerWidgetFactory(SweetJuiceWidgetFactory factory) {
        widgetFactories.put(factory.getType(), factory);
    }

    public void registerWidgetFactory(String type, SweetJuiceWidgetFactory factory) {
        widgetFactories.put(type, factory);
    }

    public void render(String jsonTree) {
        try {
            if (jsonTree == null || jsonTree.isEmpty()) {
                Log.w(TAG, "UIManager: null or empty JSON");
                return;
            }

            Log.d(TAG, "UIManager: render payload=" + jsonTree.length() + " chars");

            JSONObject rootNode;
            try {
                rootNode = new JSONObject(jsonTree);
            } catch (Exception e) {
                Log.e(TAG, "UIManager: JSON parse failed", e);
                return;
            }

            String type = rootNode.optString("type");
            if (type.isEmpty()) {
                Log.w(TAG, "UIManager: root node missing type");
                return;
            }

            renderFailed = false;

            if ("root".equals(type)) {
                JSONObject childNode = rootNode.optJSONObject("child");
                if (childNode == null) {
                    Log.w(TAG, "UIManager: root node missing child");
                    return;
                }
                View existingRoot = rootContainer.getChildAt(0);
                View newView = updateOrCreateView(existingRoot, childNode);

                if (existingRoot != newView) {
                    rootContainer.removeAllViews();
                    rootContainer.addView(newView);
                }

                applyRootBackground(rootNode);
                if (!renderFailed) {
                    Log.d(TAG, "UIManager: render OK");
                }
                return;
            }

            View existingRoot = rootContainer.getChildAt(0);
            View newView = updateOrCreateView(existingRoot, rootNode);

            if (existingRoot != newView) {
                Log.d(TAG, "UIManager: replacing root view");
                rootContainer.removeAllViews();
                rootContainer.addView(newView);
            }

            if (!renderFailed) {
                Log.d(TAG, "UIManager: render OK");
            }

        } catch (Exception e) {
            Log.e(TAG, "UIManager: render fatal", e);
        }
    }

    private void applyRootBackground(JSONObject rootNode) {
        try {
            String bg = rootNode.optString("backgroundColor", "").trim();
            if (bg.isEmpty()) {
                JSONObject style = rootNode.optJSONObject("style");
                if (style != null) {
                    bg = style.optString("backgroundColor", "").trim();
                }
            }
            if (!bg.isEmpty()) {
                int color = Color.parseColor(bg);
                rootContainer.setBackgroundColor(color);
            }
        } catch (Exception e) {
            Log.w(TAG, "UIManager: applyRootBackground failed", e);
        }
    }

    private View updateOrCreateView(View existingView, JSONObject node) {
        try {
            String id = node.optString("id", "");
            String type = node.optString("type", "");

            if (type.isEmpty()) {
                Log.w(TAG, "UIManager: node missing type, id=" + id);
                return existingView != null ? existingView : createView("text");
            }

            SweetJuiceWidgetFactory widgetFactory = widgetFactories.get(type);
            if (widgetFactory != null) {
                View widgetView = existingView;
                if (widgetView == null || !(widgetView.getTag() != null && widgetView.getTag().equals(id))) {
                    widgetView = widgetFactory.createView(context, node, rootContainer);
                    widgetView.setTag(id);
                }
                updateView(widgetView, node);
                return widgetView;
            }

            View view = existingView;
            boolean needsNewView = false;

            if (view == null) {
                needsNewView = true;
            } else if (!isViewTypeOf(view, type)) {
                needsNewView = true;
            } else if (!id.equals(view.getTag())) {
                needsNewView = true;
            }

            if (needsNewView) {
                if (existingView != null && isViewTypeOf(existingView, type)) {
                    view = existingView;
                    view.setTag(id);
                } else {
                    view = createView(type);
                    view.setTag(id);
                }
            }

            updateView(view, node);
            return view;

        } catch (Exception e) {
            Log.e(TAG, "UIManager: updateOrCreateView failed", e);
            renderFailed = true;
            return existingView != null ? existingView : createView("text");
        }
    }

    private boolean isViewTypeOf(View view, String type) {
        if (view == null) return false;
        switch (type) {
            case "column":
            case "row":
            case "button-group":
            case "segmented-button":
                return view instanceof LinearLayout && !(view instanceof MaterialCardView);
            case "card":
                return view instanceof MaterialCardView;
            case "text":
                return view instanceof TextView && !(view instanceof MaterialButton) && !(view instanceof EditText);
            case "button":
            case "text-button":
            case "outlined-button":
            case "tonal-button":
            case "elevated-button":
            case "icon-button":
                return view instanceof MaterialButton;
            case "textfield":
                return view instanceof EditText;
            case "spacer":
                return view.getClass().equals(View.class);
            case "fab":
            case "extended-fab":
                return view instanceof FloatingActionButton || view instanceof ExtendedFloatingActionButton;
            default:
                return false;
        }
    }

    private View createView(String type) {
        View v;
        try {
            switch (type) {
                case "column":
                    LinearLayout col = new LinearLayout(context);
                    col.setOrientation(LinearLayout.VERTICAL);
                    col.setBackgroundColor(Color.TRANSPARENT);
                    v = col;
                    break;
                case "row":
                    LinearLayout row = new LinearLayout(context);
                    row.setOrientation(LinearLayout.HORIZONTAL);
                    row.setBackgroundColor(Color.TRANSPARENT);
                    v = row;
                    break;
                case "card":
                    MaterialCardView card = new MaterialCardView(context);
                    card.setCardBackgroundColor(Color.WHITE);
                    card.setRadius(dpToPx(8));
                    LinearLayout cardLayout = new LinearLayout(context);
                    cardLayout.setOrientation(LinearLayout.VERTICAL);
                    cardLayout.setBackgroundColor(Color.TRANSPARENT);
                    card.addView(cardLayout);
                    v = card;
                    break;
                case "text":
                    TextView tv = new TextView(context);
                    tv.setTextColor(Color.BLACK);
                    v = tv;
                    break;
                case "button":
                    MaterialButton btn = new MaterialButton(context);
                    btn.setBackgroundColor(Color.parseColor("#6200EE"));
                    btn.setTextColor(Color.WHITE);
                    v = btn;
                    break;
                case "text-button":
                    MaterialButton textBtn = new MaterialButton(context);
                    textBtn.setBackgroundColor(Color.TRANSPARENT);
                    textBtn.setElevation(0);
                    textBtn.setStateListAnimator(null);
                    v = textBtn;
                    break;
                case "outlined-button":
                    MaterialButton outlinedBtn = new MaterialButton(context);
                    outlinedBtn.setStrokeWidth(1);
                    outlinedBtn.setStrokeColor(ColorStateList.valueOf(Color.GRAY));
                    v = outlinedBtn;
                    break;
                case "tonal-button":
                    MaterialButton tonalBtn = new MaterialButton(context);
                    tonalBtn.setBackgroundColor(Color.parseColor("#DDE1FF"));
                    v = tonalBtn;
                    break;
                case "elevated-button":
                    MaterialButton elevatedBtn = new MaterialButton(context);
                    elevatedBtn.setElevation(dpToPx(2));
                    v = elevatedBtn;
                    break;
                case "icon-button":
                    MaterialButton iconBtn = new MaterialButton(context);
                    iconBtn.setIcon(ContextCompat.getDrawable(context, android.R.drawable.ic_menu_edit));
                    iconBtn.setIconGravity(MaterialButton.ICON_GRAVITY_TEXT_START);
                    v = iconBtn;
                    break;
                case "fab":
                    v = new FloatingActionButton(context);
                    break;
                case "extended-fab":
                    ExtendedFloatingActionButton efab = new ExtendedFloatingActionButton(context);
                    v = efab;
                    break;
                case "segmented-button":
                    MaterialButtonToggleGroup group = new MaterialButtonToggleGroup(context);
                    group.setSingleSelection(true);
                    group.setSelectionRequired(true);
                    v = group;
                    break;
                case "button-group":
                    LinearLayout btnGroup = new LinearLayout(context);
                    btnGroup.setOrientation(LinearLayout.HORIZONTAL);
                    btnGroup.setBackgroundColor(Color.TRANSPARENT);
                    v = btnGroup;
                    break;
                case "textfield":
                    EditText et = new EditText(context);
                    et.setBackgroundColor(Color.LTGRAY);
                    v = et;
                    break;
                case "spacer":
                    v = new View(context);
                    v.setBackgroundColor(Color.TRANSPARENT);
                    break;
                default:
                    Log.w(TAG, "UIManager: unknown type=" + type + ", using plain View");
                    v = new View(context);
                    v.setBackgroundColor(Color.RED);
            }
        } catch (Exception e) {
            Log.e(TAG, "UIManager: createView failed for type=" + type, e);
            v = new View(context);
            v.setBackgroundColor(Color.RED);
        }

        LinearLayout.LayoutParams lp = new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.WRAP_CONTENT
        );
        v.setLayoutParams(lp);
        return v;
    }

    private void updateView(View view, JSONObject node) {
        try {
            String type = node.optString("type");
            JSONObject style = node.optJSONObject("style");
            JSONArray events = node.optJSONArray("events");
            Set<String> eventSet = new HashSet<>();
            if (events != null) {
                for (int i = 0; i < events.length(); i++) {
                    eventSet.add(events.getString(i));
                }
            }

            SweetJuiceWidgetFactory widgetFactory = widgetFactories.get(type);
            if (widgetFactory != null) {
                widgetFactory.updateView(view, node);
                setupEvents(view, node.optString("id", ""), eventSet);
                return;
            }

            if (view instanceof ViewGroup) {
                ViewGroup container;
                if (view instanceof MaterialCardView) {
                    View child = ((MaterialCardView) view).getChildAt(0);
                    if (child instanceof ViewGroup) {
                        container = (ViewGroup) child;
                    } else {
                        return;
                    }
                } else {
                    container = (ViewGroup) view;
                }

                JSONArray children = node.optJSONArray("children");
                if (children != null) {
                    int nodeCount = children.length();
                    int viewCount = container.getChildCount();

                    for (int i = 0; i < nodeCount; i++) {
                        try {
                            JSONObject childNode = children.getJSONObject(i);
                            View existingChild = (i < viewCount) ? container.getChildAt(i) : null;
                            View updatedChild = updateOrCreateView(existingChild, childNode);

                            if (existingChild == null) {
                                container.addView(updatedChild);
                            } else if (existingChild != updatedChild) {
                                container.removeViewAt(i);
                                container.addView(updatedChild, i);
                            }
                        } catch (Exception e) {
                            Log.e(TAG, "UIManager: child update failed at index " + i, e);
                        }
                    }
                    if (viewCount > nodeCount) {
                        container.removeViews(nodeCount, viewCount - nodeCount);
                    }
                } else {
                    container.removeAllViews();
                }
            }

            if (view instanceof MaterialButton) {
                MaterialButton btn = (MaterialButton) view;
                btn.setText(node.optString("text", ""));
            } else if (view instanceof FloatingActionButton) {
                FloatingActionButton fab = (FloatingActionButton) view;
                if (view instanceof ExtendedFloatingActionButton) {
                    ((ExtendedFloatingActionButton) view).setText(node.optString("text", ""));
                }
            } else if (view instanceof TextView && !(view instanceof EditText)) {
                TextView tv = (TextView) view;
                tv.setText(node.optString("value", ""));
            } else if (view instanceof EditText) {
                EditText et = (EditText) view;
                String val = node.optString("value", "");
                if (!val.equals(et.getText().toString())) {
                    et.setText(val);
                }
                et.setHint(node.optString("placeholder", ""));
            } else if ("spacer".equals(type)) {
                try {
                    float rawW = (float) node.optDouble("width", -1);
                    float rawH = (float) node.optDouble("height", -1);
                    LinearLayout.LayoutParams lp = (LinearLayout.LayoutParams) view.getLayoutParams();
                    if (lp == null) {
                        lp = new LinearLayout.LayoutParams(0, 0);
                    }
                    if (rawH >= 0) {
                        lp.height = dpToPx(rawH);
                        lp.weight = 0;
                    } else {
                        lp.height = 0;
                        lp.weight = 1.0f;
                    }
                    if (rawW >= 0) {
                        lp.width = dpToPx(rawW);
                    } else {
                        lp.width = ViewGroup.LayoutParams.MATCH_PARENT;
                    }
                    view.setLayoutParams(lp);
                } catch (Exception e) {
                    Log.e(TAG, "UIManager: spacer layout failed", e);
                }
            }

            setupEvents(view, node.optString("id", ""), eventSet);
            applyStyles(view, type, style);

        } catch (Exception e) {
            Log.e(TAG, "UIManager: updateView fatal", e);
            renderFailed = true;
        }
    }

    private void setupEvents(View view, final String id, Set<String> events) {
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
            if (view instanceof EditText) {
                EditText et = (EditText) view;
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

    private void applyStyles(View view, String type, JSONObject style) {
        try {
            if (style == null) return;

            LinearLayout.LayoutParams lp = (LinearLayout.LayoutParams) view.getLayoutParams();
            if (lp == null) {
                lp = new LinearLayout.LayoutParams(
                        ViewGroup.LayoutParams.MATCH_PARENT,
                        ViewGroup.LayoutParams.WRAP_CONTENT
                );
            }

            if (style.has("flex")) {
                try {
                    lp.weight = (float) style.optDouble("flex");
                    if (lp.weight > 0) {
                        ViewGroup parent = (ViewGroup) view.getParent();
                        boolean horizontal = parent instanceof LinearLayout &&
                                ((LinearLayout) parent).getOrientation() == LinearLayout.HORIZONTAL;
                        if (horizontal) {
                            lp.width = 0;
                            lp.height = ViewGroup.LayoutParams.MATCH_PARENT;
                        } else {
                            lp.height = 0;
                            lp.width = ViewGroup.LayoutParams.MATCH_PARENT;
                        }
                    }
                } catch (Exception e) {
                    Log.e(TAG, "UIManager: flex style failed", e);
                }
            }

            if (style.has("width")) {
                try {
                    lp.width = dpToPx((float) style.optDouble("width"));
                } catch (Exception e) {
                    Log.e(TAG, "UIManager: width style failed", e);
                }
            }
            if (style.has("height")) {
                try {
                    lp.height = dpToPx((float) style.optDouble("height"));
                } catch (Exception e) {
                    Log.e(TAG, "UIManager: height style failed", e);
                }
            }
            view.setLayoutParams(lp);

            try {
                int p = dpToPx((float) style.optDouble("padding", 0));
                int pv = dpToPx((float) style.optDouble("paddingVertical", 0));
                int ph = dpToPx((float) style.optDouble("paddingHorizontal", 0));
                int pt = pv != 0 ? pv : p;
                int pb = pv != 0 ? pv : p;
                int pl = ph != 0 ? ph : p;
                int pr = ph != 0 ? ph : p;
                if (view instanceof MaterialCardView) {
                    ((MaterialCardView) view).setContentPadding(pl, pt, pr, pb);
                } else {
                    view.setPadding(pl, pt, pr, pb);
                }
            } catch (Exception e) {
                Log.e(TAG, "UIManager: padding style failed", e);
            }

            if (style.has("backgroundColor")) {
                try {
                    int color = Color.parseColor(style.optString("backgroundColor"));
                    if (view instanceof MaterialCardView) {
                        ((MaterialCardView) view).setCardBackgroundColor(color);
                    } else {
                        view.setBackgroundColor(color);
                    }
                } catch (Exception e) {
                    Log.e(TAG, "UIManager: background color failed: " + style.optString("backgroundColor"), e);
                }
            }

            if (style.has("cornerRadius")) {
                try {
                    float radius = dpToPx((float) style.optDouble("cornerRadius"));
                    if (view instanceof MaterialCardView) {
                        ((MaterialCardView) view).setRadius(radius);
                    }
                } catch (Exception e) {
                    Log.e(TAG, "UIManager: cornerRadius failed", e);
                }
            }

            if (view instanceof LinearLayout) {
                try {
                    LinearLayout ll = (LinearLayout) view;
                    String align = style.optString("alignItems");
                    String justify = style.optString("justifyContent");
                    int gravity = 0;
                    if ("center".equals(align)) gravity |= Gravity.CENTER_HORIZONTAL;
                    if ("center".equals(justify)) gravity |= Gravity.CENTER_VERTICAL;
                    ll.setGravity(gravity);
                } catch (Exception e) {
                    Log.e(TAG, "UIManager: gravity failed", e);
                }
            }

            if (view instanceof TextView) {
                try {
                    TextView tv = (TextView) view;
                    if (style.has("fontSize")) {
                        tv.setTextSize(TypedValue.COMPLEX_UNIT_SP, (float) style.optDouble("fontSize"));
                    }
                    if (style.has("color")) {
                        tv.setTextColor(Color.parseColor(style.optString("color")));
                    }
                    if ("bold".equals(style.optString("weight"))) {
                        tv.setTypeface(null, Typeface.BOLD);
                    }
                } catch (Exception e) {
                    Log.e(TAG, "UIManager: text style failed", e);
                }
            }

        } catch (Exception e) {
            Log.e(TAG, "UIManager: applyStyles fatal", e);
            renderFailed = true;
        }
    }

    private void sendEvent(String id, String name, Object data) {
        try {
            JSONObject event = new JSONObject();
            event.put("id", id);
            event.put("name", name);
            if (data != null) event.put("data", data);
            String payload = event.toString();
            Log.d(TAG, "UIManager: sendEvent id=" + id + " name=" + name + " payload=" + payload);
            Sweetjuice.handleMessageFromFrontend("ui:event", payload);
        } catch (Exception e) {
            Log.e(TAG, "UIManager: sendEvent failed", e);
            e.printStackTrace();
        }
    }

    private int dpToPx(float dp) {
        return (int) TypedValue.applyDimension(
                TypedValue.COMPLEX_UNIT_DIP, dp, context.getResources().getDisplayMetrics()
        );
    }
}
