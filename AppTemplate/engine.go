package juiceapp

import (
	"github.com/sweet-juice/sweetjuice/app"
	"github.com/sweet-juice/sweetjuice/core"
)

// NativeCallHandler is an interface implemented by the mobile platform (Java/Obj-C)
// to handle calls originating from the Go layer.
type NativeCallHandler interface {
	OnNativeCall(method string, args string) string
}

// SetNativeCallHandler registers the platform-specific handler.
func SetNativeCallHandler(handler NativeCallHandler) {
	core.SetNativeCallHandler(handler)
}

// HandleMessageFromFrontend routes UI events and frontend calls to the Go backend.
func HandleMessageFromFrontend(methodKey string, jsonArgsPayload string) string {
	return core.HandleMessageFromFrontend(methodKey, jsonArgsPayload)
}

// ReRender triggers a UI update from the native side.
func ReRender() {
	app.ReRender()
}

// HandleNativeAction routes native plugin calls to the Go backend.
func HandleNativeAction(domain string, data string) string {
	return core.HandleNativeAction(domain, data)
}

// PollNativeEvent checks for any pending events from Go to Native.
func PollNativeEvent() string {
	return core.PollNativeEvent()
}

// RegisterPlugin records a plugin definition so the native side can instantiate it.
func RegisterPlugin(domain string, javaPkg string, className string) {
	core.RegisterPlugin(&core.PluginDefinition{
		Domain:  domain,
		JavaPkg: javaPkg,
		Class:   className,
	})
}

// GetRegisteredPlugins returns a JSON representation of all registered plugins.
func GetRegisteredPlugins() string {
	return core.GetRegisteredPluginsSnapshot()
}
