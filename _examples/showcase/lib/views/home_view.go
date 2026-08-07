package views

import (
	"myapp/lib/state"
	"github.com/sweet-juice/sweetjuice/plugins/mu3"

	"github.com/sweet-juice/sweetjuice/ui"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

type HomeView struct {
	State *state.MainAppState
}

func (v *HomeView) Render() ui.Node {
	return ui.Root(
		ui.VStack(
			mu3.TopAppBar("Feed", "menu"),
			ui.Spacer().Height(12),
			mu3.SearchBar("Search articles...", false),
			ui.Spacer().Height(16),
			mu3.Tabs([]string{"For You", "Saved", "Archived"}, nil, v.State.SelectedTab).
				OnChanged(v.State.SetSelectedTab),
			ui.Spacer().Height(16),
			v.contentForTab(),
			ui.Spacer().Height(24),
			mu3.NavigationBar([]map[string]string{
				{"label": "Home", "icon": "home"},
				{"label": "Explore", "icon": "explore"},
				{"label": "Profile", "icon": "person"},
			}),
		),
		"#FFFBFE",
	)
}

func (v *HomeView) contentForTab() ui.Node {
	switch v.State.SelectedTab {
	case "saved":
		return ui.VStack(
			mu3.Card("Design Systems", "M3 guidelines and best practices").OnClick(func() { println("saved 1") }),
			ui.Spacer().Height(12),
			mu3.Card("Go Patterns", "Idiomatic Go for mobile backends").OnClick(func() { println("saved 2") }),
		)
	case "archived":
		return ui.Text("No archived items").Style(style.Text{
			FontSize: 14,
			Color:    "#49454F",
		})
	default:
		return ui.VStack(
			mu3.Card("Getting Started", "Build your first Sweet Juice app in minutes").OnClick(func() { println("feed 1") }),
			ui.Spacer().Height(12),
			mu3.Card("Material 3 Design", "Learn about colors, typography, and shapes").OnClick(func() { println("feed 2") }),
			ui.Spacer().Height(12),
			mu3.Card("Native Components", "True native rendering with gomobile bindings").OnClick(func() { println("feed 3") }),
			ui.Spacer().Height(12),
			mu3.Card("Plugin System", "Extend your app with native modules").OnClick(func() { println("feed 4") }),
		)
	}
}
