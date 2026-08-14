package mu3

import (
	"github.com/sweet-juice/sweetjuice/core"
	"github.com/sweet-juice/sweetjuice/ui"
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

// Layout primitives
func VStack(children ...ui.Node) *ui.StackNode {
	return ui.VStack(children...)
}

func HStack(children ...ui.Node) *ui.StackNode {
	return ui.HStack(children...)
}

func Spacer() *ui.SpacerNode {
	return ui.Spacer()
}

// Text
func Text(value string) *ui.WidgetNode {
	return ui.Widget("mu3:text", map[string]interface{}{
		"value": value,
	})
}

// Buttons
func Button(text string) *ui.WidgetNode {
	return ui.Widget("mu3:button", map[string]interface{}{
		"text": text,
	})
}

func TextButton(text string) *ui.WidgetNode {
	return ui.Widget("mu3:text-button", map[string]interface{}{
		"text": text,
	})
}

func OutlinedButton(text string) *ui.WidgetNode {
	return ui.Widget("mu3:outlined-button", map[string]interface{}{
		"text": text,
	})
}

func TonalButton(text string) *ui.WidgetNode {
	return ui.Widget("mu3:tonal-button", map[string]interface{}{
		"text": text,
	})
}

func ElevatedButton(text string) *ui.WidgetNode {
	return ui.Widget("mu3:elevated-button", map[string]interface{}{
		"text": text,
	})
}

func IconButton(name string) *ui.WidgetNode {
	return ui.Widget("mu3:icon-button", map[string]interface{}{
		"name": name,
	})
}

func SegmentedButton(options []string, selected string) *ui.WidgetNode {
	return ui.Widget("mu3:segmented-button", map[string]interface{}{
		"options":  options,
		"selected": selected,
	})
}

func ButtonGroup(children ...ui.Node) *ui.WidgetNode {
	childMaps := make([]map[string]interface{}, len(children))
	for i, c := range children {
		childMaps[i], _ = c.Serialize()
	}
	return ui.Widget("mu3:button-group", map[string]interface{}{
		"children": childMaps,
	})
}

// FABs
func FAB() *ui.WidgetNode {
	return ui.Widget("mu3:fab", map[string]interface{}{})
}

func ExtendedFAB(text string) *ui.WidgetNode {
	return ui.Widget("mu3:extended-fab", map[string]interface{}{
		"text": text,
	})
}

// Input
func TextField(placeholder string) *ui.WidgetNode {
	return ui.Widget("mu3:textfield", map[string]interface{}{
		"placeholder": placeholder,
	})
}

// Media
func Image(src string) *ui.WidgetNode {
	return ui.Widget("mu3:image", map[string]interface{}{
		"src": src,
	})
}

func Video(src string) *ui.WidgetNode {
	return ui.Widget("mu3:video", map[string]interface{}{
		"src": src,
	})
}

// Containers
func Card(title string, subtitle string) *ui.WidgetNode {
	return ui.Widget("mu3:card", map[string]interface{}{
		"title":    title,
		"subtitle": subtitle,
	})
}

func Box(children ...ui.Node) *ui.WidgetNode {
	return &ui.WidgetNode{
		BaseNode: ui.BaseNode{Type: "mu3:box", ID: ui.GenID()},
		Props:    make(map[string]interface{}),
		Children: children,
	}
}

func Root(child ui.Node, backgroundColor string) *ui.RootNode {
	return ui.Root(child, backgroundColor)
}

// Icons
func IconOutlined(name string) *ui.WidgetNode {
	return ui.Widget("mu3:icon-outlined", map[string]interface{}{
		"name": name,
	})
}

// Material 3 specific components
func TopAppBar(title string, navigationIcon string) *ui.WidgetNode {
	return ui.Widget("mu3:top-app-bar", map[string]interface{}{
		"title":          title,
		"navigationIcon": navigationIcon,
	})
}

func BottomAppBar() *ui.WidgetNode {
	return ui.Widget("mu3:bottom-app-bar", map[string]interface{}{})
}

func NavigationBar(destinations []map[string]string) *ui.WidgetNode {
	return ui.Widget("mu3:nav-bar", map[string]interface{}{
		"destinations": destinations,
	})
}

func NavigationRail(destinations []map[string]string) *ui.WidgetNode {
	return ui.Widget("mu3:nav-rail", map[string]interface{}{
		"destinations": destinations,
	})
}

func SearchBar(hint string, showDocked bool) *ui.WidgetNode {
	return ui.Widget("mu3:search-bar", map[string]interface{}{
		"hint":       hint,
		"showDocked": showDocked,
	})
}

func Tabs(primary []string, secondary []string, selected string) *ui.WidgetNode {
	return ui.Widget("mu3:tabs", map[string]interface{}{
		"primary":   primary,
		"secondary": secondary,
		"selected":  selected,
	})
}

func Toolbar(items []map[string]string, orientation string) *ui.WidgetNode {
	return ui.Widget("mu3:toolbar", map[string]interface{}{
		"items":       items,
		"orientation": orientation,
	})
}

func (p *Plugin) ShowCard(cardJSON string) string {
	return core.CallNativePlatform("mu3:showCard", cardJSON)
}

func Dialog(title, message, buttonText string) *ui.WidgetNode {
	return ui.Widget("mu3:dialog", map[string]interface{}{
		"title":      title,
		"message":    message,
		"buttonText": buttonText,
	})
}
