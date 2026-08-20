package app

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/sweet-juice/sweetjuice/core"
	"github.com/sweet-juice/sweetjuice/ui"
)

var rootComponent ui.Component

func Run(c ui.Component) {
	rootComponent = c

	a := core.NewApplication(core.Options{
		Name: "Sweet Juice App",
		OnStart: func(app *core.Application) error {
			return ReRender()
		},
	})

	core.SetUIEventDispatcher(func(id string, event string, data interface{}) error {
		ui.DispatchEvent(id, event, data)
		return nil
	})

	if err := a.Run(); err != nil {
		fmt.Printf("App failed to start: %v\n", err)
	}
}

func GetEventBus() *core.EventBus {
	app := core.GetGlobalApp()
	if app == nil {
		return nil
	}
	return app.Events
}

func EmitEvent(name string, data interface{}) {
	bus := GetEventBus()
	if bus != nil {
		bus.Emit(name, data)
	}
}

func GetGlobalApp() *core.Application {
	return core.GetGlobalApp()
}

func sanitizePayload(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(val))
		for k, vv := range val {
			result[k] = sanitizePayload(vv)
		}
		return result
	case []map[string]interface{}:
		result := make([]map[string]interface{}, len(val))
		for i, vv := range val {
			result[i] = sanitizePayload(vv).(map[string]interface{})
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, vv := range val {
			result[i] = sanitizePayload(vv)
		}
		return result
	case []string:
		return val
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Func {
			return nil
		}
		return v
	}
}

func ReRender() error {
	fmt.Println("Go: app.ReRender called")
	if rootComponent == nil {
		fmt.Println("Go: Error - rootComponent is nil")
		return fmt.Errorf("no root component set")
	}

	node := rootComponent.Render()
	return RenderNode(node)
}

func RenderNode(node ui.Node) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	tree, err := node.Serialize()
	if err != nil {
		fmt.Printf("Go: Serialization error: %v\n", err)
		return err
	}

	tree = sanitizePayload(tree).(map[string]interface{})

	payload, err := json.Marshal(tree)
	if err != nil {
		fmt.Printf("Go: JSON Marshal error after sanitize: %v\n", err)
		return err
	}

	if len(payload) == 0 {
		fmt.Println("Go: Warning - empty JSON payload generated")
		return fmt.Errorf("empty JSON payload")
	}

	fmt.Printf("Go: Dispatching UI render (%d bytes)\n", len(payload))
	core.CallNativePlatform("ui:render", string(payload))
	return nil
}

func ShowOverlay(child ui.Node) string {
	node := ui.Overlay(child)
	RenderNode(node)
	return node.BaseNode.ID
}

func DismissOverlay(id string) {
	node := ui.DismissOverlay(id)
	RenderNode(node)
}
