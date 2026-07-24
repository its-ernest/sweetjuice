package app

import (
	"encoding/json"
	"fmt"
	"github.com/sweet-juice/sweetjuice/core"
	"github.com/sweet-juice/sweetjuice/ui"
)

var rootComponent ui.Component

// Run initializes the application with a root component.
func Run(c ui.Component) {
	rootComponent = c

	// We create a core application under the hood
	a := core.NewApplication(core.Options{
		Name: "Sweet Juice App",
		OnStart: func(app *core.Application) error {
			return ReRender()
		},
	})

	// Set global event dispatcher for UI events
	core.SetUIEventDispatcher(func(id string, event string, data map[string]interface{}) error {
		ui.DispatchEvent(id, event, data)
		return nil
	})

	if err := a.Run(); err != nil {
		fmt.Printf("App failed to start: %v\n", err)
	}
}

// ReRender triggers a UI update by rendering the root component again.
func ReRender() error {
	if rootComponent == nil {
		return fmt.Errorf("no root component set")
	}

	node := rootComponent.Render()
	tree, err := node.Serialize()
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(tree)
	core.CallNativePlatform("ui:render", string(payload))
	return nil
}
