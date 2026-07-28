package com.sweetjuice.pkg.mu3;

import android.content.Context;
import android.view.View;
import android.view.ViewGroup;
import com.sweetjuice.plugin.SweetJuicePlugin;
import com.sweetjuice.plugin.SweetJuiceWidgetFactory;
import org.json.JSONObject;
import org.json.JSONException;

public class Mu3Plugin implements SweetJuicePlugin {

    @Override
    public String getDomain() {
        return "mu3";
    }

    @Override
    public void onAttach(Context context) {
        // no-op
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if ("widget:register".equals(action)) {
            return registerWidgets();
        }
        return "{}";
    }

    @Override
    public SweetJuiceWidgetFactory[] getWidgetFactories() {
        return new SweetJuiceWidgetFactory[] {
            new Mu3CardWidgetFactory()
        };
    }

    private String registerWidgets() {
        try {
            JSONObject result = new JSONObject();
            result.put("registered", true);
            return result.toString();
        } catch (JSONException e) {
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }
}
