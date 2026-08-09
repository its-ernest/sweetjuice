package com.sweetjuice.pkg.sms;

import android.content.Context;
import android.database.Cursor;
import android.net.Uri;
import android.os.Build;
import android.provider.Telephony;
import android.util.Log;
import com.sweetjuice.plugin.SweetJuicePlugin;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class SmsPlugin implements SweetJuicePlugin {

    private static final String TAG = "SmsPlugin";
    private Context mContext;

    @Override
    public String getDomain() { return "sms"; }

    @Override
    public void onAttach(Context context) {
        this.mContext = context;
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        try {
            if ("getInbox".equals(action)) {
                return getSmsFolder(Telephony.Sms.Inbox.CONTENT_URI, "inbox");
            }
            if ("getSent".equals(action)) {
                return getSmsFolder(Telephony.Sms.Sent.CONTENT_URI, "sent");
            }
            if ("getDrafts".equals(action)) {
                return getSmsFolder(Telephony.Sms.Draft.CONTENT_URI, "draft");
            }
            if ("getAll".equals(action)) {
                return getAllSms();
            }
            return new JSONObject().put("error", "Unknown action").toString();
        } catch (JSONException e) {
            return errorJson(e.getMessage());
        }
    }

    private String getSmsFolder(Uri uri, String folder) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.KITKAT) {
            String defaultSmsPackage = Telephony.Sms.getDefaultSmsPackage(mContext);
            if (defaultSmsPackage != null && !defaultSmsPackage.equals(mContext.getPackageName())) {
                return errorJson("Not default SMS app");
            }
        }

        Cursor cursor = null;
        try {
            cursor = mContext.getContentResolver().query(uri, null, null, null, "date DESC LIMIT 100");
            if (cursor == null) {
                return folderJson(folder, 0, new JSONArray());
            }

            JSONArray messages = new JSONArray();
            int count = 0;

            while (cursor.moveToNext()) {
                try {
                    JSONObject msg = new JSONObject();
                    msg.put("id", cursor.getLong(cursor.getColumnIndexOrThrow("_id")));
                    msg.put("address", cursor.getString(cursor.getColumnIndexOrThrow("address")));
                    msg.put("body", cursor.getString(cursor.getColumnIndexOrThrow("body")));
                    msg.put("type", folder);
                    msg.put("timestamp", cursor.getLong(cursor.getColumnIndexOrThrow("date")));
                    msg.put("read", cursor.getInt(cursor.getColumnIndexOrThrow("read")) == 1);
                    messages.put(msg);
                    count++;
                } catch (JSONException e) {
                    Log.w(TAG, "Failed to parse SMS row", e);
                }
            }

            return folderJson(folder, count, messages);
        } catch (Exception e) {
            return errorJson(e.getMessage());
        } finally {
            if (cursor != null) {
                cursor.close();
            }
        }
    }

    private String getAllSms() {
        JSONArray allMessages = new JSONArray();
        int totalCount = 0;

        String[] folders = {"inbox", "sent", "draft"};
        Uri[] uris = {Telephony.Sms.Inbox.CONTENT_URI, Telephony.Sms.Sent.CONTENT_URI, Telephony.Sms.Draft.CONTENT_URI};

        for (int i = 0; i < folders.length; i++) {
            String result = getSmsFolder(uris[i], folders[i]);
            try {
                JSONObject obj = new JSONObject(result);
                if (obj.has("messages")) {
                    JSONArray msgs = obj.getJSONArray("messages");
                    for (int j = 0; j < msgs.length(); j++) {
                        allMessages.put(msgs.getJSONObject(j));
                    }
                    totalCount += obj.optInt("count", 0);
                }
            } catch (JSONException e) {
                Log.w(TAG, "Failed to merge SMS folder: " + folders[i], e);
            }
        }

        try {
            JSONObject result = new JSONObject();
            result.put("folder", "all");
            result.put("count", totalCount);
            result.put("messages", allMessages);
            return result.toString();
        } catch (JSONException e) {
            return errorJson(e.getMessage());
        }
    }

    private String folderJson(String folder, int count, JSONArray messages) {
        try {
            JSONObject obj = new JSONObject();
            obj.put("folder", folder);
            obj.put("count", count);
            obj.put("messages", messages);
            return obj.toString();
        } catch (JSONException e) {
            return errorJson(e.getMessage());
        }
    }

    private String okJson(String status) {
        try {
            return new JSONObject().put("status", status).toString();
        } catch (JSONException e) {
            return "{\"status\":\"" + status + "\"}";
        }
    }

    private String errorJson(String message) {
        try {
            return new JSONObject().put("error", message).toString();
        } catch (JSONException e) {
            return "{\"error\":\"" + message + "\"}";
        }
    }

    @Override public void onRequestPermissionsResult(int rc, String[] p, int[] g) {}
    @Override public void onActivityResult(int r, int rc, android.content.Intent d) {}
    @Override public void onNewIntent(android.content.Intent intent) {}
}
