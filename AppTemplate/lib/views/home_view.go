package views

import (
	"fmt"
	"time"

	"sweetjuice/lib/state"

	"github.com/sweet-juice/sweetjuice/plugins/mu3"
	"github.com/sweet-juice/sweetjuice/plugins/permission"
	"github.com/sweet-juice/sweetjuice/plugins/special"
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
		mu3.Button("Request Call & SMS").OnClick(func() {
			fmt.Println("HomeView: requesting call and SMS permissions")
			v.requestCallAndSMS()
		}),
		ui.Spacer().Height(12),
		mu3.Button("Request Location").OnClick(func() {
			fmt.Println("HomeView: requesting location permissions")
			v.requestLocation()
		}),
		ui.Spacer().Height(12),
		mu3.Button("Request Background Location").OnClick(func() {
			fmt.Println("HomeView: requesting background location permission")
			v.requestBackgroundLocation()
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

	return ui.Root(root, "#FFFBFE")
}

func (v *HomeView) requestCallAndSMS() {
	v.State.SetWaiting("Requesting call and SMS permissions...")
	go func() {
		_, err := permission.RequestMultiple([]string{
			"android.permission.READ_CALL_LOG",
			"android.permission.READ_SMS",
		})
		if err != nil {
			fmt.Println("HomeView: RequestMultiple error:", err)
		}
		v.State.ClearWaiting()
	}()
}

func (v *HomeView) requestLocation() {
	v.State.SetWaiting("Requesting location permissions...")
	go func() {
		_, err := permission.RequestMultiple([]string{
			"android.permission.ACCESS_FINE_LOCATION",
			"android.permission.ACCESS_COARSE_LOCATION",
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
