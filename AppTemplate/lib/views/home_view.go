package views

import (
	"myapp/lib/components"
	"myapp/lib/state"
	"github.com/sweet-juice/sweetjuice/ui"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

type HomeView struct {
	State *state.MainAppState
}

func (v *HomeView) Render() ui.Node {
	return ui.VStack(
		components.Header("SweetJuice"),

		ui.Spacer().Height(40),

		components.UserInput(v.State.User),
	).
		Style(style.View{
			Flex:            1,
			Padding:         32,
			JustifyContent: style.Center,
			AlignItems:      style.Center,
			BackgroundColor: "#FFFFFF",
		})
}
