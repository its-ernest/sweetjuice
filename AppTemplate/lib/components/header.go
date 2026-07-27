package components

import (
	"github.com/sweet-juice/sweetjuice/ui"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

func Header(title string) ui.Node {
	return ui.Text(title).
		Style(style.Text{
			FontSize: 28,
			Weight:   style.Bold,
			Color:    "#1A1A1A",
		})
}
