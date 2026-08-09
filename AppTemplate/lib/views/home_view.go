package views

import (
	"encoding/json"
	"fmt"
	"time"

	"sweetjuice/lib/state"

	"github.com/sweet-juice/sweetjuice/plugins/calls"
	"github.com/sweet-juice/sweetjuice/plugins/mu3"
	"github.com/sweet-juice/sweetjuice/plugins/permission"
	"github.com/sweet-juice/sweetjuice/plugins/sms"
	"github.com/sweet-juice/sweetjuice/plugins/special"
	"github.com/sweet-juice/sweetjuice/plugins/workmanager"
	"github.com/sweet-juice/sweetjuice/ui"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

type HomeView struct {
	State *state.MainAppState
}

func (v *HomeView) Render() ui.Node {
	root := ui.VStack(
		ui.Spacer().Height(24),
		mu3.TopAppBar("Sweet Juice", ""),
		ui.Spacer().Height(24),
		ui.Text("Grant permissions to continue").Style(style.Text{
			FontSize: 18,
			Color:    "#1C1B1F",
		}),
		ui.Spacer().Height(24),
	)

	root = ui.VStack(
		root,
		mu3.Button("Request Permissions").OnClick(func() {
			fmt.Println("HomeView: requesting runtime permissions")
			v.requestRuntimePermissions()
		}),
		ui.Spacer().Height(12),
		mu3.Button("Request Background Location").OnClick(func() {
			fmt.Println("HomeView: requesting background location permission")
			v.requestBackgroundLocation()
		}),
		ui.Spacer().Height(12),
		mu3.Button("Exempt App").OnClick(func() {
			fmt.Println("HomeView: requesting battery exemption")
			v.requestBatteryExemption()
		}),
		ui.Spacer().Height(12),
		mu3.Button("Get All Calls").OnClick(func() {
			fmt.Println("HomeView: fetching all calls")
			v.getAllCalls()
		}),
		ui.Spacer().Height(12),
		mu3.Button("Get All SMS").OnClick(func() {
			fmt.Println("HomeView: fetching all SMS")
			v.getAllSms()
		}),
		ui.Spacer().Height(12),
		mu3.Button("Start Task").OnClick(func() {
			fmt.Println("HomeView: starting notification task")
			v.startNotificationTask()
		}),
		ui.Spacer().Height(24),
	)

	if v.State.Waiting {
		root = ui.VStack(
			root,
			ui.Text(v.State.WaitingMsg).Style(style.Text{
				FontSize: 14,
				Color:    "#7A757E",
			}),
			ui.Spacer().Height(16),
		)
	}

	if v.State.CallLogMsg != "" {
		root = ui.VStack(
			root,
			ui.Text("Calls Result:").Style(style.Text{
				FontSize: 16,
				Color:    "#1C1B1F",
			}),
			ui.Text(v.State.CallLogMsg).Style(style.Text{
				FontSize: 12,
				Color:    "#4A4458",
			}),
			ui.Spacer().Height(12),
		)
	}

	if v.State.SmsMsg != "" {
		root = ui.VStack(
			root,
			ui.Text("SMS Result:").Style(style.Text{
				FontSize: 16,
				Color:    "#1C1B1F",
			}),
			ui.Text(v.State.SmsMsg).Style(style.Text{
				FontSize: 12,
				Color:    "#4A4458",
			}),
			ui.Spacer().Height(12),
		)
	}

	return ui.Root(root, "#FFFBFE")
}

func (v *HomeView) requestRuntimePermissions() {
	v.State.SetWaiting("Requesting permissions...")
	go func() {
		_, err := permission.RequestMultiple([]string{
			"android.permission.READ_CALL_LOG",
			"android.permission.READ_SMS",
			"android.permission.ACCESS_FINE_LOCATION",
			"android.permission.ACCESS_COARSE_LOCATION",
			"android.permission.POST_NOTIFICATIONS",
		})
		if err != nil {
			fmt.Println("HomeView: RequestMultiple error:", err)
		}
		v.State.ClearWaiting()
	}()
}

func (v *HomeView) requestBackgroundLocation() {
	v.State.SetWaiting("Requesting background location permission...")
	go func() {
		status, err := permission.Check("android.permission.ACCESS_FINE_LOCATION")
		if err != nil {
			fmt.Println("HomeView: Check error:", err)
		}
		if status != "granted" {
			fmt.Println("HomeView: foreground location not granted, opening app settings")
			v.State.ClearWaiting()
			return
		}

		_, err = permission.Request("android.permission.ACCESS_BACKGROUND_LOCATION")
		if err != nil {
			fmt.Println("HomeView: Request error:", err)
		}

		time.Sleep(2 * time.Second)
		status, _ = permission.Check("android.permission.ACCESS_BACKGROUND_LOCATION")
		if status != "granted" {
			fmt.Println("HomeView: background location not granted after request, opening app settings")
			_, _ = special.RequestAppSettings()
		}
		v.State.ClearWaiting()
	}()
}

func (v *HomeView) requestBatteryExemption() {
	v.State.SetWaiting("Requesting battery exemption...")
	go func() {
		_, err := special.RequestBatteryExemption()
		if err != nil {
			fmt.Println("HomeView: BatteryExemption error:", err)
		}
		v.State.ClearWaiting()
	}()
}

func (v *HomeView) getAllCalls() {
	v.State.SetWaiting("Fetching all calls...")
	go func() {
		plugin := calls.NewPlugin()
		log, err := plugin.GetAll()
		if err != nil {
			fmt.Println("HomeView: GetAll calls error:", err)
			v.State.SetCallLogResult(calls.CallLog{}, fmt.Sprintf("Error: %v", err))
			return
		}

		b, _ := json.MarshalIndent(log, "", "  ")
		summary := fmt.Sprintf("Count: %d\n%s", log.Count, string(b))
		fmt.Println("HomeView: Call Log Result:\n", summary)
		v.State.SetCallLogResult(log, summary)
	}()
}

func (v *HomeView) getAllSms() {
	v.State.SetWaiting("Fetching all SMS...")
	go func() {
		plugin := sms.NewPlugin()
		folder, err := plugin.GetAll()
		if err != nil {
			fmt.Println("HomeView: GetAll SMS error:", err)
			v.State.SetSmsResult(sms.SmsFolder{}, fmt.Sprintf("Error: %v", err))
			return
		}

		b, _ := json.MarshalIndent(folder, "", "  ")
		summary := fmt.Sprintf("Folder: %s\nCount: %d\n%s", folder.Folder, folder.Count, string(b))
		fmt.Println("HomeView: SMS Result:\n", summary)
		v.State.SetSmsResult(folder, summary)
	}()
}

func (v *HomeView) startNotificationTask() {
	const taskKey = "sweetjuice_notification_task"
	const intervalMinutes = 15

	_, _ = workmanager.NewPlugin().EnqueuePeriodic(taskKey, intervalMinutes, nil, false)
}
