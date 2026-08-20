package juiceapp

import (
	"sweetjuice/lib/state"
	"sweetjuice/lib/views"

	"github.com/sweet-juice/sweetjuice/app"
)

func StartApplication() string {
	registerPluginDefinitions()

	mainState := state.NewMainAppState()
	root := &views.HomeView{State: mainState}

	// app.Run starts the core runtime and triggers the first render.
	app.Run(root)

	// initPlugins registers Go-side native methods on the current app instance.
	initPlugins()

	return `{"status":"started"}`
}
