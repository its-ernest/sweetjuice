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
type textNode struct {
	BaseNode
	Value string
}

func Text(value string) *textNode {
	return &textNode{
		BaseNode: BaseNode{Type: "text", ID: genID()},
		Value:    value,
	}
}

func (n *textNode) ID(id string) *textNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *textNode) Style(s style.Text) *textNode {
	n.BaseNode.Style = s
	return n
}

func (n *textNode) OnClick(h func()) *textNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *textNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"value":  n.Value,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

// Button Node
type buttonNode struct {
	BaseNode
	Text string
}

func Button(text string) *buttonNode {
	return &buttonNode{
		BaseNode: BaseNode{Type: "button", ID: genID()},
		Text:     text,
	}
}

func (n *buttonNode) ID(id string) *buttonNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *buttonNode) Style(s style.Button) *buttonNode {
	n.BaseNode.Style = s
	return n
}

func (n *buttonNode) OnClick(h func()) *buttonNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *buttonNode) OnLongPress(h func()) *buttonNode {
	n.register("long_press", func(interface{}) { h() })
	return n
}

func (n *buttonNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"text":   n.Text,
		"style":  n.BaseNode.Style,
		"events": n.Events,
	}, nil
}

// TextField Node
type textFieldNode struct {
	BaseNode
	Placeholder string
	Value       string
}

func TextField(placeholder string) *textFieldNode {
	return &textFieldNode{
		BaseNode:    BaseNode{Type: "textfield", ID: genID()},
		Placeholder: placeholder,
	}
}

func (n *textFieldNode) ID(id string) *textFieldNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *textFieldNode) WithValue(v string) *textFieldNode {
	n.Value = v
	return n
}

func (n *textFieldNode) OnChanged(h func(string)) *textFieldNode {
	n.register("changed", func(args interface{}) {
		if s, ok := args.(string); ok {
			n.Value = s
			h(s)
		}
	})
	return n
}

func (n *textFieldNode) Style(s style.View) *textFieldNode {
	n.BaseNode.Style = s
	return n
}

func (n *textFieldNode) Serialize() (map[string]interface{}, error) {
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
type stackNode struct {
	BaseNode
	Children []Node
}

func VStack(children ...Node) *stackNode {
	return &stackNode{
		BaseNode: BaseNode{Type: "column", ID: genID()},
		Children: children,
	}
}

func HStack(children ...Node) *stackNode {
	return &stackNode{
		BaseNode: BaseNode{Type: "row", ID: genID()},
		Children: children,
	}
}

func (n *stackNode) ID(id string) *stackNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *stackNode) Style(s style.View) *stackNode {
	n.BaseNode.Style = s
	return n
}

func (n *stackNode) OnClick(h func()) *stackNode {
	n.register("click", func(interface{}) { h() })
	return n
}

func (n *stackNode) Serialize() (map[string]interface{}, error) {
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
type cardNode struct {
	stackNode
}

func Card(children ...Node) *cardNode {
	return &cardNode{
		stackNode: stackNode{
			BaseNode: BaseNode{Type: "card", ID: genID()},
			Children: children,
		},
	}
}

func (n *cardNode) ID(id string) *cardNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *cardNode) Style(s style.View) *cardNode {
	n.BaseNode.Style = s
	return n
}

// Spacer Node
type spacerNode struct {
	BaseNode
	width  float64
	height float64
}

func Spacer() *spacerNode {
	return &spacerNode{
		BaseNode: BaseNode{Type: "spacer", ID: genID()},
	}
}

func (n *spacerNode) ID(id string) *spacerNode {
	n.BaseNode.SetID(id)
	return n
}

func (n *spacerNode) Width(w float64) *spacerNode {
	n.width = w
	return n
}

func (n *spacerNode) Height(h float64) *spacerNode {
	n.height = h
	return n
}

func (n *spacerNode) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":   n.Type,
		"id":     n.BaseNode.ID,
		"width":  n.width,
		"height": n.height,
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
