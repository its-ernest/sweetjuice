package views

import (
	"fmt"
	"myapp/lib/state"
	"github.com/sweet-juice/sweetjuice/plugins/mu3"
	"github.com/sweet-juice/sweetjuice/plugins/permission"

	"github.com/sweet-juice/sweetjuice/ui"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

type HomeView struct {
	State *state.MainAppState
}

func (v *HomeView) Render() ui.Node {
	name := v.State.User.Name
	greeting := "Please enter your name"
	if name != "" {
		greeting = fmt.Sprintf("Hello, %s!", name)
	}

	return ui.Root(
		ui.VStack(
			ui.Text(greeting).
				Style(style.Text{
					FontSize: 22,
					Weight:   style.Bold,
					Color:    "#FFFFFF",
				}),

			ui.Spacer().Height(24),

			mu3.Card("Sweet Juice", "Material 3 Components Demo").OnClick(func() {
				println("mu3 card tapped")
			}),

			ui.Spacer().Height(24),

			ui.TextField("Your Name").
				ID("user_name_input").
				WithValue(name).
				OnChanged(v.State.User.SetName).
				Style(style.View{
					BackgroundColor: "#F2F2F7",
					Padding:         16,
					CornerRadius:    12,
				}),

			ui.Spacer().Height(24),

			ui.Text("Action Buttons").
				Style(style.Text{
					FontSize: 16,
					Weight:   style.Bold,
					Color:    "#FFFFFF",
				}),

			ui.Spacer().Height(12),

			ui.VStack(
				ui.Button("Filled Button").
					OnClick(func() { println("filled clicked") }).
					Style(style.Button{
						BackgroundColor:   "#6750A4",
						Color:             "#FFFFFF",
						CornerRadius:      20,
						PaddingHorizontal: 24,
						PaddingVertical:   12,
					}),

				ui.OutlinedButton("Outlined Button").
					OnClick(func() { println("outlined clicked") }).
					Style(style.OutlinedButton{
						Button: style.Button{
							Color:             "#6750A4",
							PaddingHorizontal: 24,
							PaddingVertical:   12,
						},
						StrokeWidth: 1,
						StrokeColor: "#6750A4",
					}),

				ui.TextButton("Text Button").
					OnClick(func() { println("text clicked") }).
					Style(style.TextButton{
						Button: style.Button{
							Color:             "#FFC107",
							PaddingHorizontal: 24,
							PaddingVertical:   12,
						},
					}),

				ui.TonalButton("Tonal Button").
					OnClick(func() { println("tonal clicked") }),

				ui.ElevatedButton("Elevated Button").
					OnClick(func() { println("elevated clicked") }),
			),

			ui.Spacer().Height(24),

			ui.Text("Icon Button & FAB").
				Style(style.Text{
					FontSize: 16,
					Weight:   style.Bold,
					Color:    "#FFFFFF",
				}),

			ui.Spacer().Height(12),

			ui.HStack(
				ui.IconButton("favorite").
					OnClick(func() { println("icon clicked") }),
				ui.Spacer().Width(16),
				mu3.StandardFAB().OnClick(func() { println("fab clicked") }),
			),

			ui.Spacer().Height(24),

			ui.Text("Segmented Control").
				Style(style.Text{
					FontSize: 16,
					Weight:   style.Bold,
					Color:    "#FFFFFF",
				}),

			ui.Spacer().Height(12),

			ui.SegmentedButton([]string{"All", "Active", "Completed"}, "All").
				OnChanged(func(selected string) {
					fmt.Printf("selected: %s\n", selected)
				}),

			ui.Spacer().Height(24),

			ui.Button("Enable Notifications").
				OnClick(func() {
					status, err := permission.Request("android.permission.POST_NOTIFICATIONS")
					fmt.Printf("Notification permission status: %s, err: %v\n", status, err)
				}).
				Style(style.Button{
					BackgroundColor:   "#FF9500",
					CornerRadius:      20,
					PaddingHorizontal: 24,
					PaddingVertical:   12,
				}),

			ui.Spacer().Height(24),

			ui.Button("Say Hello").
				OnClick(func() {
					fmt.Printf("User says: %s\n", v.State.User.Name)
				}).
				Style(style.Button{
					BackgroundColor:   "#34C759",
					CornerRadius:      20,
					PaddingHorizontal: 24,
					PaddingVertical:   12,
				}),
		).
			Style(style.View{
				Padding:        32,
				JustifyContent: style.Center,
				AlignItems:     style.Center,
			}),
		"#1A1A1A",
	)
}
