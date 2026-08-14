package app

import (
	"encoding/json"

	"github.com/sweet-juice/sweetjuice/core"
)

// AppPlugin handles app lifecycle events.
type AppPlugin struct {
	app *core.Application
}

// NewAppPlugin creates a new AppPlugin.
func NewAppPlugin() *AppPlugin {
	return &AppPlugin{}
}

// Init registers the native callback for app lifecycle events.
func (p *AppPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("app:resumed", func(args []json.RawMessage) (interface{}, error) {
		var payload map[string]interface{}
		if len(args) > 0 {
			json.Unmarshal(args[0], &payload)
		}

		app.Events.Emit("app:resumed", payload)
		return map[string]string{"status": "received"}, nil
	})

	return nil
}

// Name returns the plugin name.
func (p *AppPlugin) Name() string {
	return "app"
}

// Domain returns the plugin domain.
func (p *AppPlugin) Domain() string {
	return "app"
}
