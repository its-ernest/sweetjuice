package com.sweetjuice.pkg.mu3;

import android.content.Context;
import com.sweetjuice.plugin.SweetJuicePlugin;
import com.sweetjuice.plugin.SweetJuiceWidgetFactory;

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
        return "{}";
    }

    @Override
    public SweetJuiceWidgetFactory[] getWidgetFactories() {
        return new SweetJuiceWidgetFactory[] {
            new Mu3TextWidgetFactory(),
            new Mu3ButtonWidgetFactory(),
            new Mu3TextButtonWidgetFactory(),
            new Mu3OutlinedButtonWidgetFactory(),
            new Mu3TonalButtonWidgetFactory(),
            new Mu3ElevatedButtonWidgetFactory(),
            new Mu3IconButtonWidgetFactory(),
            new Mu3FabWidgetFactory(),
            new Mu3ExtendedFabWidgetFactory(),
            new Mu3SegmentedButtonWidgetFactory(),
            new Mu3ButtonGroupWidgetFactory(),
            new Mu3TextFieldWidgetFactory(),
            new Mu3ImageWidgetFactory(),
            new Mu3VideoWidgetFactory(),
            new Mu3CardWidgetFactory(),
            new Mu3SpacerWidgetFactory(),
            new Mu3IconOutlinedWidgetFactory(),
            new Mu3TopAppBarWidgetFactory(),
            new Mu3BottomAppBarWidgetFactory(),
            new Mu3NavigationBarWidgetFactory(),
            new Mu3NavigationRailWidgetFactory(),
            new Mu3SearchBarWidgetFactory(),
            new Mu3TabsWidgetFactory(),
            new Mu3ToolbarWidgetFactory()
        };
    }
}
