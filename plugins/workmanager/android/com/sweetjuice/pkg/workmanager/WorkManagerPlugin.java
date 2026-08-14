package com.sweetjuice.pkg.workmanager;

import android.content.Context;
import android.content.Intent;
import androidx.work.Constraints;
import androidx.work.Data;
import androidx.work.NetworkType;
import androidx.work.OneTimeWorkRequest;
import androidx.work.PeriodicWorkRequest;
import androidx.work.WorkInfo;
import androidx.work.WorkManager;

import com.sweetjuice.plugin.SweetJuicePlugin;

import org.json.JSONException;
import org.json.JSONObject;

import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

/**
 * WorkManagerPlugin (WorkManager) allows Go to schedule background tasks.
 */
public class WorkManagerPlugin implements SweetJuicePlugin {
    private Context mContext;

    @Override
    public String getDomain() { return "workmanager"; }

    @Override
    public void onAttach(Context context) { 
        this.mContext = context; 
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        try {
            JSONObject args = new JSONObject(jsonArgsPayload);
            String taskKey = args.optString("task_key", "default_task");
            
            Constraints constraints = parseConstraints(args.optJSONObject("constraints"));

            if ("enqueueOneTime".equals(action)) {
                Data inputData = new Data.Builder().putString("task_key", taskKey).build();

                long initialDelayMinutes = args.optLong("initial_delay_minutes", 0);
                long initialDelaySeconds = args.optLong("initial_delay_seconds", 0);
                boolean replaceExisting = args.optBoolean("replace_existing", false);

                if (replaceExisting) {
                    WorkManager.getInstance(mContext).cancelAllWorkByTag(taskKey);
                }

                OneTimeWorkRequest.Builder builder = new OneTimeWorkRequest.Builder(SweetJuiceBackgroundWorker.class)
                        .setInputData(inputData)
                        .setConstraints(constraints)
                        .addTag(taskKey);

                if (initialDelaySeconds > 0) {
                    builder.setInitialDelay(initialDelaySeconds, TimeUnit.SECONDS);
                } else if (initialDelayMinutes > 0) {
                    builder.setInitialDelay(initialDelayMinutes, TimeUnit.MINUTES);
                }

                OneTimeWorkRequest request = builder.build();

                WorkManager.getInstance(mContext).enqueue(request);
                return "{\"status\":\"enqueued\",\"id\":\"" + request.getId().toString() + "\"}";
            }

            if ("enqueuePeriodic".equals(action)) {
                long intervalMinutes = args.optLong("interval_minutes", 15);
                boolean replaceExisting = args.optBoolean("replace_existing", false);
                boolean runImmediate = args.optBoolean("run_immediate", false);
                Data inputData = new Data.Builder().putString("task_key", taskKey).build();

                if (replaceExisting) {
                    WorkManager.getInstance(mContext).cancelAllWorkByTag(taskKey);
                }

                PeriodicWorkRequest request = new PeriodicWorkRequest.Builder(
                        SweetJuiceBackgroundWorker.class, intervalMinutes, TimeUnit.MINUTES)
                        .setInputData(inputData)
                        .setConstraints(constraints)
                        .addTag(taskKey)
                        .build();

                WorkManager.getInstance(mContext).enqueue(request);
                String response = "{\"status\":\"periodic_enqueued\",\"id\":\"" + request.getId().toString() + "\"}";

                if (runImmediate) {
                    Data immediateInput = new Data.Builder().putString("task_key", taskKey).build();
                    OneTimeWorkRequest immediateRequest = new OneTimeWorkRequest.Builder(SweetJuiceBackgroundWorker.class)
                            .setInputData(immediateInput)
                            .setConstraints(constraints)
                            .addTag(taskKey)
                            .build();
                    WorkManager.getInstance(mContext).enqueue(immediateRequest);
                    response = response.replace("}", ",\"immediate_id\":\"" + immediateRequest.getId().toString() + "\"}");
                }

                return response;
            }

            if ("isEnqueued".equals(action)) {
                try {
                    List<WorkInfo> workInfos = WorkManager.getInstance(mContext).getWorkInfosByTag(taskKey).get();
                    boolean enqueued = false;
                    for (WorkInfo info : workInfos) {
                        if (info.getState() == WorkInfo.State.ENQUEUED || 
                            info.getState() == WorkInfo.State.RUNNING || 
                            info.getState() == WorkInfo.State.BLOCKED) {
                            enqueued = true;
                            break;
                        }
                    }
                    return "{\"enqueued\":" + enqueued + "}";
                } catch (ExecutionException | InterruptedException e) {
                    return "{\"error\":\"" + e.getMessage() + "\"}";
                }
            }

            if ("cancelAll".equals(action)) {
                WorkManager.getInstance(mContext).cancelAllWork();
                return "{\"status\":\"cancelled_all\"}";
            }

        } catch (JSONException e) {
            return "{\"error\":\"Invalid JSON payload\"}";
        }

        return "{\"error\":\"Unknown action\"}";
    }

    private Constraints parseConstraints(JSONObject json) {
        Constraints.Builder builder = new Constraints.Builder();
        if (json == null) return builder.build();

        String net = json.optString("network_type", "NOT_REQUIRED");
        switch (net) {
            case "CONNECTED": builder.setRequiredNetworkType(NetworkType.CONNECTED); break;
            case "UNMETERED": builder.setRequiredNetworkType(NetworkType.UNMETERED); break;
            case "NOT_ROAMING": builder.setRequiredNetworkType(NetworkType.NOT_ROAMING); break;
            default: builder.setRequiredNetworkType(NetworkType.NOT_REQUIRED);
        }

        builder.setRequiresCharging(json.optBoolean("requires_charging", false));
        builder.setRequiresDeviceIdle(json.optBoolean("requires_device_idle", false));
        builder.setRequiresBatteryNotLow(json.optBoolean("requires_battery_not_low", false));
        builder.setRequiresStorageNotLow(json.optBoolean("requires_storage_not_low", false));

        return builder.build();
    }

    @Override public void onRequestPermissionsResult(int rc, String[] p, int[] g) {}
    @Override public void onActivityResult(int r, int rc, Intent d) {}
    @Override public void onNewIntent(Intent intent) {}
}
