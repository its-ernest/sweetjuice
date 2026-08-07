package sweetjuice

import (
	"github.com/sweet-juice/sweetjuice/core"
)

// NativeCallHandler is an interface that Java can implement to handle calls from Go.
// Defining it here ensures gomobile generates it in the 'sweetjuice' Java package.
type NativeCallHandler interface {
	OnNativeCall(method string, args string) string
}

// SetNativeCallHandler registers the Java-side handler for Go-to-Native calls.
func SetNativeCallHandler(handler NativeCallHandler) {
	core.SetNativeCallHandler(handler)
}

// Below functions are called from Java to handle messages/events from the frontend or to perform native actions.

// HandleMessageFromFrontend processes messages sent from the JavaScript frontend.
func HandleMessageFromFrontend(methodKey string, jsonArgsPayload string) string {
	return core.HandleMessageFromFrontend(methodKey, jsonArgsPayload)
}

// HandleNativeAction processes calls from Go to Java and returns the result back to Go.
func HandleNativeAction(methodKey string, jsonArgsPayload string) string {
	return core.HandleNativeAction(methodKey, jsonArgsPayload)
}

// PollNativeEvent allows Go to check for any events sent from Java and retrieve them as a JSON string.
func PollNativeEvent() string {
	return core.NewMobileBridge().PollNativeEvent()
}
