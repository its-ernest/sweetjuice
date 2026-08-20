package juiceapp

import (
	"fmt"

	"github.com/sweet-juice/sweetjuice/app"
	apppkg "github.com/sweet-juice/sweetjuice/plugins/app"
	"github.com/sweet-juice/sweetjuice/plugins/broadcast"
	"github.com/sweet-juice/sweetjuice/plugins/mu3"
	"github.com/sweet-juice/sweetjuice/plugins/system"
)

// registerPluginDefinitions registers plugin metadata with Java.
func registerPluginDefinitions() {
	RegisterPlugin("app", "com.sweetjuice.pkg.app", "AppPlugin")
	RegisterPlugin("mu3", "com.sweetjuice.pkg.mu3", "Mu3Plugin")
	RegisterPlugin("broadcast", "com.sweetjuice.pkg.broadcast", "BroadcastPlugin")
	RegisterPlugin("system", "com.sweetjuice.pkg.system", "SystemPlugin")
}

// initPlugins initializes all plugins with the app instance.
func initPlugins() {
	a := app.GetGlobalApp()
	if a == nil {
		fmt.Println("initPlugins: global app not ready")
		return
	}

	apppkg.NewAppPlugin().Init(a)
	mu3.New().Init(a)
	broadcast.NewPlugin().Init(a)
	system.NewPlugin().Init(a)
}
