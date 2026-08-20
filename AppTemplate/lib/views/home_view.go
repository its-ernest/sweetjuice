package views

import (
	"sweetjuice/lib/state"

	"github.com/sweet-juice/sweetjuice/app"
	"github.com/sweet-juice/sweetjuice/plugins/mu3"
	"github.com/sweet-juice/sweetjuice/ui"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

type HomeView struct {
	State *state.MainAppState
}

func (v *HomeView) Render() ui.Node {
	root := ui.VStack(
		ui.Spacer().Height(48),
		ui.Text("Sweet Juice Tests").Style(style.Text{FontSize: 24, Weight: style.Bold}),
		ui.Spacer().Height(24),

		mu3.Button("Show Native Dialog").OnClick(func() {
			v.showTestDialog()
		}),

		ui.Spacer().Height(16),

		mu3.Button("Show Custom Overlay").OnClick(func() {
			v.showTestOverlay()
		}),
	).Style(style.View{
		Padding:    24,
		AlignItems: style.Center,
	})

	return ui.Root(root, "#FFFBFE")
}

func (v *HomeView) showTestDialog() {
	dialog := ui.Dialog("Test Dialog", "This is a native Material 3 alert dialog.").
		WithConfirm("Got it")

	app.RenderNode(dialog)
}

func (v *HomeView) showTestOverlay() {
	var overlayID string

	content := mu3.Box(
		ui.Text("Custom Overlay").Style(style.Text{FontSize: 20, Weight: style.Bold}),
		ui.Spacer().Height(12),
		ui.Text("This is a flexible Box rendered in the overlay layer.").Style(style.Text{FontSize: 16}),
		ui.Spacer().Height(24),
		mu3.Button("Close Overlay").OnClick(func() {
			app.DismissOverlay(overlayID)
		}),
	).Style(style.View{
		Padding:         24,
		BackgroundColor: "#FFFFFF",
		CornerRadius:    28,
		Width:           300,
	})

	overlayID = app.ShowOverlay(ui.VStack(
		ui.Spacer(),
		ui.HStack(ui.Spacer(), content, ui.Spacer()),
		ui.Spacer(),
	))
}
