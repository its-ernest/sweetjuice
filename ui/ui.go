package ui

import (
	"crypto/rand"
	"fmt"
	"github.com/sweet-juice/sweetjuice/ui/style"
)

type Node interface {
	Serialize() (map[string]interface{}, error)
}

type Component interface {
	Render() Node
}

type BaseNode struct {
	Type   string
	ID     string
	Style  interface{}
	Events []string
}

func genID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("n_%x", b)
}

func (b *BaseNode) SetID(id string) {
	b.ID = id
}

func (b *BaseNode) register(event string, handler func(interface{})) {
	b.Events = append(b.Events, event)
	RegisterEvent(b.ID, event, handler)
}

// Text Node
type TextNode struct {
	BaseNode
	Value string
}

func Text(value string) *TextNode {
	return &TextNode{
		BaseNode: BaseNode{Type: "mu3:text", ID: genID()},
		Value:    value,
	}
}

func (n *TextNode) ID(id string) *TextNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *TextNode) Style(s style.Text) *TextNode {
	n.BaseNode.Style = s
	return n
}

func (n *TextNode) OnClick(h func()) *TextNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *TextNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"value":  n.Value,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

// Button Node
type ButtonNode struct {
	BaseNode
	Text string
}

func Button(text string) *ButtonNode {
	return &ButtonNode{
		BaseNode: BaseNode{Type: "mu3:button", ID: genID()},
		Text:     text,
	}
}

func (n *ButtonNode) ID(id string) *ButtonNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *ButtonNode) Style(s style.Button) *ButtonNode {
	n.BaseNode.Style = s
	return n
}

func (n *ButtonNode) OnClick(h func()) *ButtonNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *ButtonNode) OnLongPress(h func()) *ButtonNode {
	n.register("long_press", func(interface{}) { h() })
	return n
}

func (n *ButtonNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"text":   n.Text,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

// TextField Node
type TextFieldNode struct {
	BaseNode
	Placeholder string
	Value       string
}

func TextField(placeholder string) *TextFieldNode {
	return &TextFieldNode{
		BaseNode:    BaseNode{Type: "mu3:textfield", ID: genID()},
		Placeholder: placeholder,
	}
}

func (n *TextFieldNode) ID(id string) *TextFieldNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *TextFieldNode) WithValue(v string) *TextFieldNode {
	n.Value = v
	return n
}

func (n *TextFieldNode) OnChanged(h func(string)) *TextFieldNode {
	n.register("changed", func(args interface{}) {
		if s, ok := args.(string); ok {
			n.Value = s
			h(s)
		}
	})
	return n
}

func (n *TextFieldNode) Style(s style.View) *TextFieldNode {
	n.BaseNode.Style = s
	return n
}

func (n *TextFieldNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":        n.Type,
		"id":          n.BaseNode.ID,
		"placeholder": n.Placeholder,
		"value":       n.Value,
		"style":       n.BaseNode.Style,
		"events":      n.Events,
	}, nil
}

// Stack Node (VStack, HStack)
type StackNode struct {
	BaseNode
	Children []Node
}

func VStack(children ...Node) *StackNode {
	return &StackNode{
		BaseNode: BaseNode{Type: "column", ID: genID()},
		Children: children,
	}
}

func HStack(children ...Node) *StackNode {
	return &StackNode{
		BaseNode: BaseNode{Type: "row", ID: genID()},
		Children: children,
	}
}

func (n *StackNode) ID(id string) *StackNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *StackNode) Style(s style.View) *StackNode {
	n.BaseNode.Style = s
	return n
}

func (n *StackNode) OnClick(h func()) *StackNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *StackNode) Serialize() (map[string]interface{}, error) {
	children := make([]map[string]interface{}, len(n.Children))
	for i, c := range n.Children {
		children[i], _ = c.Serialize()
	}
	return map[string]interface{}{
		"type":     n.Type,
		"id":       n.BaseNode.ID,
		"children": children,
		"style":    n.BaseNode.Style,
		"events":   n.Events,
	}, nil
}

// Card Node
type CardNode struct {
	StackNode
}

func Card(children ...Node) *CardNode {
	return &CardNode{
		StackNode: StackNode{
			BaseNode: BaseNode{Type: "mu3:card", ID: genID()},
			Children: children,
		},
	}
}

func (n *CardNode) ID(id string) *CardNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *CardNode) Style(s style.View) *CardNode {
	n.BaseNode.Style = s
	return n
}

// Spacer Node
type SpacerNode struct {
	BaseNode
	width  float64
	height float64
}

func Spacer() *SpacerNode {
	return &SpacerNode{
		BaseNode: BaseNode{Type: "spacer", ID: genID()},
	}
}

func (n *SpacerNode) ID(id string) *SpacerNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *SpacerNode) Width(w float64) *SpacerNode {
	n.width = w
	return n
}

func (n *SpacerNode) Height(h float64) *SpacerNode {
	n.height = h
	return n
}

func (n *SpacerNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"width":  n.width,
		"height": n.height,
	}, nil
}

// Root Node — wraps the entire app tree and carries root-level metadata like background color.
type RootNode struct {
	BaseNode
	Child           Node
	BackgroundColor string
}

func Root(child Node, backgroundColor string) *RootNode {
	return &RootNode{
		BaseNode:        BaseNode{Type: "root", ID: genID()},
		Child:           child,
		BackgroundColor: backgroundColor,
	}
}

func (n *RootNode) Serialize() (map[string]interface{}, error) {
	childMap, err := n.Child.Serialize()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"type":             n.Type,
		"id":               n.BaseNode.ID,
		"child":            childMap,
		"events":           n.Events,
		"style":            n.BaseNode.Style,
		"backgroundColor":  n.BackgroundColor,
	}, nil
}

// Widget Node (for third-party widget plugins)
type WidgetNode struct {
	BaseNode
	Props map[string]interface{}
}

func Widget(widgetType string, props map[string]interface{}) *WidgetNode {
	if props == nil {
		props = make(map[string]interface{})
	}
	return &WidgetNode{
		BaseNode: BaseNode{Type: widgetType, ID: genID()},
		Props:    props,
	}
}

func (n *WidgetNode) ID(id string) *WidgetNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *WidgetNode) Prop(key string, value interface{}) *WidgetNode {
	n.Props[key] = value
	return n
}

func (n *WidgetNode) Style(s interface{}) *WidgetNode {
	n.BaseNode.Style = s
	return n
}

func (n *WidgetNode) OnClick(h func()) *WidgetNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *WidgetNode) OnLongPress(h func()) *WidgetNode {
	n.register("long_press", func(interface{}) { h() })
	return n
}

func (n *WidgetNode) OnChanged(h func(string)) *WidgetNode {
	n.register("changed", func(args interface{}) {
		if s, ok := args.(string); ok {
			h(s)
		}
	})
	return n
}

func (n *WidgetNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"props":  n.Props,
		"events": n.Events,
	}, nil
}

func TextButton(text string) *TextButtonNode {
	return &TextButtonNode{
		BaseNode: BaseNode{Type: "mu3:text-button", ID: genID()},
		Text:     text,
	}
}

type TextButtonNode struct {
	BaseNode
	Text string
}

func (n *TextButtonNode) ID(id string) *TextButtonNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *TextButtonNode) Style(s style.TextButton) *TextButtonNode {
	n.BaseNode.Style = s
	return n
}

func (n *TextButtonNode) OnClick(h func()) *TextButtonNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *TextButtonNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"text":   n.Text,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

func OutlinedButton(text string) *OutlinedButtonNode {
	return &OutlinedButtonNode{
		BaseNode: BaseNode{Type: "mu3:outlined-button", ID: genID()},
		Text:     text,
	}
}

type OutlinedButtonNode struct {
	BaseNode
	Text string
}

func (n *OutlinedButtonNode) ID(id string) *OutlinedButtonNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *OutlinedButtonNode) Style(s style.OutlinedButton) *OutlinedButtonNode {
	n.BaseNode.Style = s
	return n
}

func (n *OutlinedButtonNode) OnClick(h func()) *OutlinedButtonNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *OutlinedButtonNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"text":   n.Text,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

func TonalButton(text string) *TonalButtonNode {
	return &TonalButtonNode{
		BaseNode: BaseNode{Type: "mu3:tonal-button", ID: genID()},
		Text:     text,
	}
}

type TonalButtonNode struct {
	BaseNode
	Text string
}

func (n *TonalButtonNode) ID(id string) *TonalButtonNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *TonalButtonNode) Style(s style.TonalButton) *TonalButtonNode {
	n.BaseNode.Style = s
	return n
}

func (n *TonalButtonNode) OnClick(h func()) *TonalButtonNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *TonalButtonNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"text":   n.Text,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

func ElevatedButton(text string) *ElevatedButtonNode {
	return &ElevatedButtonNode{
		BaseNode: BaseNode{Type: "mu3:elevated-button", ID: genID()},
		Text:     text,
	}
}

type ElevatedButtonNode struct {
	BaseNode
	Text string
}

func (n *ElevatedButtonNode) ID(id string) *ElevatedButtonNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *ElevatedButtonNode) Style(s style.ElevatedButton) *ElevatedButtonNode {
	n.BaseNode.Style = s
	return n
}

func (n *ElevatedButtonNode) OnClick(h func()) *ElevatedButtonNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *ElevatedButtonNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"text":   n.Text,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

func IconButton(name string) *IconButtonNode {
	return &IconButtonNode{
		BaseNode: BaseNode{Type: "mu3:icon-button", ID: genID()},
		Name:     name,
	}
}

type IconButtonNode struct {
	BaseNode
	Name string
}

func (n *IconButtonNode) ID(id string) *IconButtonNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *IconButtonNode) Style(s style.IconButton) *IconButtonNode {
	n.BaseNode.Style = s
	return n
}

func (n *IconButtonNode) OnClick(h func()) *IconButtonNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *IconButtonNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"name":   n.Name,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

func FAB() *FabNode {
	return &FabNode{
		BaseNode: BaseNode{Type: "mu3:fab", ID: genID()},
	}
}

func ExtendedFAB(text string) *FabNode {
	return &FabNode{
		BaseNode: BaseNode{Type: "mu3:extended-fab", ID: genID()},
		Text:     text,
	}
}

type FabNode struct {
	BaseNode
	Text string
}

func (n *FabNode) ID(id string) *FabNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *FabNode) Style(s style.FabStyle) *FabNode {
	n.BaseNode.Style = s
	return n
}

func (n *FabNode) OnClick(h func()) *FabNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *FabNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"text":   n.Text,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

func SegmentedButton(options []string, selected string) *SegmentedButtonNode {
	return &SegmentedButtonNode{
		BaseNode: BaseNode{Type: "mu3:segmented-button", ID: genID()},
		Options:  options,
		Selected: selected,
	}
}

type SegmentedButtonNode struct {
	BaseNode
	Options  []string
	Selected string
}

func (n *SegmentedButtonNode) ID(id string) *SegmentedButtonNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *SegmentedButtonNode) Style(s style.SegmentedButton) *SegmentedButtonNode {
	n.BaseNode.Style = s
	return n
}

func (n *SegmentedButtonNode) OnChanged(h func(string)) *SegmentedButtonNode {
	n.register("changed", func(args interface{}) {
		if s, ok := args.(string); ok {
			n.Selected = s
			h(s)
		}
	})
	return n
}

func (n *SegmentedButtonNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":       n.Type,
		"id":         n.BaseNode.ID,
		"options":    n.Options,
		"selected":   n.Selected,
		"style":      n.BaseNode.Style,
		"events":     n.Events,
	}, nil
}

func ButtonGroup(children ...Node) *ButtonGroupNode {
	return &ButtonGroupNode{
		BaseNode: BaseNode{Type: "mu3:button-group", ID: genID()},
		Children: children,
	}
}

type ButtonGroupNode struct {
	BaseNode
	Children []Node
}

func (n *ButtonGroupNode) ID(id string) *ButtonGroupNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *ButtonGroupNode) Style(s style.ButtonGroup) *ButtonGroupNode {
	n.BaseNode.Style = s
	return n
}

func (n *ButtonGroupNode) Serialize() (map[string]interface{}, error) {
	children := make([]map[string]interface{}, len(n.Children))
	for i, c := range n.Children {
		children[i], _ = c.Serialize()
	}
	return map[string]interface{}{
		"type":     n.Type,
		"id":       n.BaseNode.ID,
		"children": children,
		"style":    n.BaseNode.Style,
		"events":   n.Events,
	}, nil
}

// Image Node
type ImageNode struct {
	BaseNode
	Src string
}

func Image(src string) *ImageNode {
	return &ImageNode{
		BaseNode: BaseNode{Type: "mu3:image", ID: genID()},
		Src:      src,
	}
}

func (n *ImageNode) ID(id string) *ImageNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *ImageNode) Style(s style.Image) *ImageNode {
	n.BaseNode.Style = s
	return n
}

func (n *ImageNode) OnClick(h func()) *ImageNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *ImageNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"src":    n.Src,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

// Video Node
type VideoNode struct {
	BaseNode
	Src      string
	Autoplay bool
	Loop     bool
	Muted    bool
	Controls bool
}

func Video(src string) *VideoNode {
	return &VideoNode{
		BaseNode: BaseNode{Type: "mu3:video", ID: genID()},
		Src:      src,
	}
}

func (n *VideoNode) ID(id string) *VideoNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *VideoNode) Style(s style.Video) *VideoNode {
	n.BaseNode.Style = s
	return n
}

func (n *VideoNode) OnClick(h func()) *VideoNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *VideoNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":      n.Type,
		"id":        n.BaseNode.ID,
		"src":       n.Src,
		"autoplay":  n.Autoplay,
		"loop":      n.Loop,
		"muted":     n.Muted,
		"controls":  n.Controls,
		"style":     n.BaseNode.Style,
		"events":    n.Events,
	}, nil
}

// Dialog Node - renders as a native AlertDialog overlay
type DialogNode struct {
	BaseNode
	Title      string
	Message    string
	ButtonText string
}

func Dialog(title, message, buttonText string) *DialogNode {
	return &DialogNode{
		BaseNode:   BaseNode{Type: "ui:dialog", ID: genID()},
		Title:      title,
		Message:    message,
		ButtonText: buttonText,
	}
}

func (n *DialogNode) ID(id string) *DialogNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *DialogNode) OnConfirm(h func(interface{})) *DialogNode {
	n.register("confirm", func(data interface{}) { h(data) })
	return n
}

func (n *DialogNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":       n.Type,
		"id":         n.BaseNode.ID,
		"title":      n.Title,
		"message":    n.Message,
		"buttonText": n.ButtonText,
		"events":     n.Events,
	}, nil
}

// Event Registry
var eventRegistry = make(map[string]func(interface{}))

func RegisterEvent(id, event string, handler func(interface{})) {
	eventRegistry[fmt.Sprintf("%s:%s", id, event)] = handler
}

func DispatchEvent(id, event string, args interface{}) {
	key := fmt.Sprintf("%s:%s", id, event)
	if h, ok := eventRegistry[key]; ok {
		fmt.Printf("Go: DispatchEvent hit key=%s\n", key)
		h(args)
	} else {
		fmt.Printf("Go: DispatchEvent MISS key=%s registryKeys=%v\n", key, getRegistryKeys())
	}
}

func getRegistryKeys() []string {
	keys := make([]string, 0, len(eventRegistry))
	for k := range eventRegistry {
		keys = append(keys, k)
	}
	return keys
}
