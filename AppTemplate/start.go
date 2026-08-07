package juiceapp

import (
	"sweetjuice/lib/state"
	"sweetjuice/lib/views"

	"github.com/sweet-juice/sweetjuice/app"
)

// StartApplication is the bootstrap function called by the native Android/iOS layer.
func StartApplication() string {
	mainState := state.NewMainAppState()

	root := &views.HomeView{
		State: mainState,
	}

	registerPluginDefinitions()

	app.Run(root)

	initPlugins()

	return `{"status":"started"}`
}
