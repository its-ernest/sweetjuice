package mu3

import (
	"github.com/sweet-juice/sweetjuice/core"
	"github.com/sweet-juice/sweetjuice/ui"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

type Plugin struct {
	app *core.Application
}

func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

func Card(title string, subtitle string) *ui.WidgetNode {
	return ui.Widget("mu3:card", map[string]interface{}{
		"title":    title,
		"subtitle": subtitle,
	})
}

func FilledButton(text string) *ui.ButtonNode {
	return ui.Button(text).Style(style.Button{BackgroundColor: "#6750A4", Color: "#FFFFFF"})
}

func OutlinedButton(text string) *ui.OutlinedButtonNode {
	return ui.OutlinedButton(text).Style(style.OutlinedButton{Button: style.Button{Color: "#6750A4"}})
}

func TextButton(text string) *ui.TextButtonNode {
	return ui.TextButton(text).Style(style.TextButton{Button: style.Button{Color: "#6750A4"}})
}

func TonalButton(text string) *ui.TonalButtonNode {
	return ui.TonalButton(text)
}

func ElevatedButton(text string) *ui.ElevatedButtonNode {
	return ui.ElevatedButton(text)
}

func IconButton(name string) *ui.IconButtonNode {
	return ui.IconButton(name)
}

func StandardFAB() *ui.FabNode {
	return ui.FAB()
}

func ExtendedFAB(text string) *ui.FabNode {
	return ui.ExtendedFAB(text)
}

func SegmentedButton(options []string, selected string) *ui.SegmentedButtonNode {
	return ui.SegmentedButton(options, selected)
}

func ButtonGroup(children ...ui.Node) *ui.ButtonGroupNode {
	return ui.ButtonGroup(children...)
}

func (p *Plugin) ShowCard(cardJSON string) string {
	return core.CallNativePlatform("mu3:showCard", cardJSON)
}
