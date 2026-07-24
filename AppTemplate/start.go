package sweetjuice

import (
	"helloworld/lib/state"
	"helloworld/lib/views"
	"github.com/sweet-juice/sweetjuice/app"
)

// StartApplication is the bootstrap function called by the native Android/iOS layer.
func StartApplication() string {
	// 1. Initialize our main state container
	mainState := state.NewMainAppState()

	// 2. Initialize the root view with the state
	root := &views.HomeView{
		State: mainState,
	}

	// 3. Launch the SweetJuice application engine
	app.Run(root)

	return `{"status":"started"}`
}

// ReRender allows the native side to trigger a UI refresh when the activity lifecycle changes.
func ReRender() {
	app.ReRender()
}
