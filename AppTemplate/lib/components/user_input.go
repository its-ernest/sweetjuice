package components

import (
	"fmt"
	"helloworld/lib/state"
	"github.com/sweet-juice/sweetjuice/ui"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

func UserInput(s *state.UserState) ui.Node {
	greeting := "Please enter your name"
	if s.Name != "" {
		greeting = fmt.Sprintf("Hello, %s!", s.Name)
	}

	return ui.VStack(
		ui.Text(greeting).
			Style(style.Text{
				FontSize: 24,
				Weight:   style.Bold,
				Color:    "#007AFF",
			}),

		ui.Spacer().Height(24),

		ui.TextField("Your Name").
			OnChanged(s.SetName).
			Style(style.View{
				BackgroundColor: "#F2F2F7",
				Padding:         16,
				CornerRadius:    12,
			}),

		ui.Spacer().Height(24),

		ui.Button("Say Hello").
			OnClick(func() {
				fmt.Printf("User says: %s\n", s.Name)
			}).
			Style(style.Button{
				BackgroundColor: "#34C759",
				CornerRadius:    8,
				PaddingHorizontal: 24,
				PaddingVertical: 12,
			}),
	).Style(style.View{
		AlignItems: style.Center,
	})
}
