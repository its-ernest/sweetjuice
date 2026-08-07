package com.sweetjuice.pkg.permissions;

import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;

import com.sweetjuice.plugin.SweetJuicePlugin;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayDeque;
import java.util.Queue;

import sweetjuice.Sweetjuice;

/**
 * PermissionsPlugin handles system runtime permissions.
 * It is background-aware for status checks but requires an active Activity for request prompts.
 */
public class PermissionsPlugin implements SweetJuicePlugin {
    private Context mContext;
    private AppCompatActivity mActivity;
    private static final int PERMISSION_REQ_CODE = 9911;
    private static final String TAG = "PermissionsPlugin";
    private final Queue<String> pendingPermissions = new ArrayDeque<>();

    @Override
    public String getDomain() { return "permissions"; }

    @Override
    public void onAttach(Context context) {
        this.mContext = context;
        if (context instanceof AppCompatActivity) {
            this.mActivity = (AppCompatActivity) context;
        }
        Log.d(TAG, "onAttach activity=" + mActivity);
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        Log.d(TAG, "handleAction action=" + action + " payload=" + jsonArgsPayload + " activity=" + mActivity);
        if ("check".equals(action)) {
            String perm = parsePermissionFromJson(jsonArgsPayload);
            int result = ContextCompat.checkSelfPermission(mContext, perm);
            return result == PackageManager.PERMISSION_GRANTED ? "{\"status\":\"granted\"}" : "{\"status\":\"denied\"}";
        }

        if ("request".equals(action)) {
            if (mActivity == null) {
                return "{\"error\":\"No active UI to request permissions\"}";
            }
            String perm = parsePermissionFromJson(jsonArgsPayload);
            Log.d(TAG, "requestSingle permission=" + perm);
            ActivityCompat.requestPermissions(mActivity, new String[]{perm}, PERMISSION_REQ_CODE);
            return "{\"status\":\"requested\"}";
        }

        if ("requestMultiple".equals(action)) {
            if (mActivity == null) {
                return "{\"error\":\"No active UI to request permissions\"}";
            }
            try {
                JSONObject obj = new JSONObject(jsonArgsPayload);
                JSONArray perms = obj.getJSONArray("permissions");
                pendingPermissions.clear();
                for (int i = 0; i < perms.length(); i++) {
                    pendingPermissions.offer(perms.getString(i));
                }
                Log.d(TAG, "requestMultiple queued=" + pendingPermissions.size());
                requestNextPermission();
                return "{\"status\":\"requesting\"}";
            } catch (JSONException e) {
                Log.e(TAG, "requestMultiple parse failed", e);
                return "{\"error\":\"" + e.getMessage() + "\"}";
            }
        }

        return "{\"error\":\"Unknown action\"}";
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {
        Log.d(TAG, "onRequestPermissionsResult code=" + requestCode + " count=" + permissions.length);
        if (requestCode == PERMISSION_REQ_CODE && grantResults.length > 0) {
            boolean granted = grantResults[0] == PackageManager.PERMISSION_GRANTED;
            String permission = permissions[0];

            JSONObject result = new JSONObject();
            try {
                result.put("permission", permission);
                result.put("granted", granted);
                String payload = "[" + result.toString() + "]";
                Sweetjuice.handleNativeAction("permissions:result", payload);
            } catch (JSONException e) {
                Log.e(TAG, "Error creating result JSON", e);
            }

            if (!pendingPermissions.isEmpty()) {
                Log.d(TAG, "requesting next permission, remaining=" + pendingPermissions.size());
                requestNextPermission();
            } else {
                Log.d(TAG, "no more permissions in queue");
            }
        }
    }

    private void requestNextPermission() {
        String perm = pendingPermissions.poll();
        Log.d(TAG, "requestNextPermission perm=" + perm + " remaining=" + pendingPermissions.size() + " activity=" + mActivity);
        if (perm != null && mActivity != null) {
            ActivityCompat.requestPermissions(mActivity, new String[]{perm}, PERMISSION_REQ_CODE);
        }
    }

    @Override public void onActivityResult(int r, int rc, Intent d) {}
    @Override public void onNewIntent(android.content.Intent intent) {}

    private String parsePermissionFromJson(String json) {
        try {
            JSONObject obj = new JSONObject(json);
            return obj.optString("permission", android.Manifest.permission.CAMERA);
        } catch (JSONException e) {
            return android.Manifest.permission.CAMERA;
        }
    }
}
