package com.sweetjuice.app;

import android.content.Intent;
import android.os.Bundle;
import android.view.ViewGroup;
import android.widget.LinearLayout;
import android.widget.ScrollView;
import android.util.Log;
import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;
import com.sweetjuice.plugin.SweetJuicePlugin;
import sweetjuice.Sweetjuice;

public class SweetJuiceActivity extends AppCompatActivity {

    private UIManager mUIManager;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        // Register this activity as active in application lifecycle
        ((SweetJuiceApplication) getApplication()).setActiveActivity(this);

        // Attach this activity context to all plugins
        SweetJuiceApplication app = (SweetJuiceApplication) getApplication();
        for (SweetJuicePlugin plugin : app.getPlugins().values()) {
            plugin.onAttach(this);
        }

        // Initialize Native UI Container
        ScrollView scrollView = new ScrollView(this);
        scrollView.setFillViewport(true);
        LinearLayout rootLayout = new LinearLayout(this);
        rootLayout.setOrientation(LinearLayout.VERTICAL);
        scrollView.addView(rootLayout);
        
        setContentView(scrollView);

        mUIManager = new UIManager(this, rootLayout);

        // Force a re-render now that the activity is active and UI is attached
        Sweetjuice.reRender();

        // Handle initial intent for Deep Linking
        handleIntent(getIntent());
    }

    public void renderUI(final String json) {
        Log.d("SweetJuice", "Activity.renderUI called");
        runOnUiThread(() -> mUIManager.render(json));
    }

    @Override
    protected void onNewIntent(Intent intent) {
        super.onNewIntent(intent);
        setIntent(intent);
        handleIntent(intent);
    }

    private void handleIntent(Intent intent) {
        if (intent == null) return;
        SweetJuiceApplication app = (SweetJuiceApplication) getApplication();
        for (SweetJuicePlugin plugin : app.getPlugins().values()) {
            plugin.onNewIntent(intent);
        }
    }

    @Override
    protected void onPause() {
        super.onPause();
    }

    @Override
    protected void onResume() {
        super.onResume();
        ((SweetJuiceApplication) getApplication()).setActiveActivity(this);
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        SweetJuiceApplication app = (SweetJuiceApplication) getApplication();
        for (SweetJuicePlugin plugin : app.getPlugins().values()) {
            plugin.onRequestPermissionsResult(requestCode, permissions, grantResults);
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        SweetJuiceApplication app = (SweetJuiceApplication) getApplication();
        for (SweetJuicePlugin plugin : app.getPlugins().values()) {
            plugin.onActivityResult(requestCode, resultCode, data);
        }
    }

    @Override
    protected void onDestroy() {
        if (((SweetJuiceApplication) getApplication()).getActiveActivity() == this) {
            ((SweetJuiceApplication) getApplication()).setActiveActivity(null);
        }
        super.onDestroy();
    }
}
