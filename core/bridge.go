package core

import (
	"encoding/json"
	"fmt"
)

var globalAppInstance *Application

// NativeCallHandler is an interface implemented by the mobile platform (Java/Obj-C)
// to handle calls originating from the Go layer.
type NativeCallHandler interface {
	OnNativeCall(method string, args string) string
}

var globalNativeHandler NativeCallHandler

// SetNativeCallHandler registers the platform-specific handler.
func SetNativeCallHandler(handler NativeCallHandler) {
	globalNativeHandler = handler
}

// CallNativePlatform calls the registered native handler from Go.
func CallNativePlatform(method string, args string) string {
	if globalNativeHandler == nil {
		return `{"error": "No native handler registered"}`
	}
	return globalNativeHandler.OnNativeCall(method, args)
}

func SetGlobalApp(app *Application) {
	globalAppInstance = app
}

var uiEventDispatcher func(id string, event string, data interface{}) error

// SetUIEventDispatcher registers the global dispatcher function to handle native UI events in Go.
func SetUIEventDispatcher(dispatcher func(id string, event string, data interface{}) error) {
	uiEventDispatcher = dispatcher
}

// MobileBridge is a dedicated wrapper struct that gobind can parse into a Java Class.
type MobileBridge struct{}

func NewMobileBridge() *MobileBridge {
	return &MobileBridge{}
}

// CallGoBackend provides a class-mapped method for the JNI layer to invoke.
func (b *MobileBridge) CallGoBackend(methodKey string, jsonArgsPayload string) string {
	fmt.Printf("Go: CallGoBackend method=%s payload=%s\n", methodKey, jsonArgsPayload)

	if globalAppInstance == nil {
		return `{"error": "Application core runtime context not active"}`
	}

	if methodKey == "ui:event" {
		if uiEventDispatcher == nil {
			fmt.Println("Go: uiEventDispatcher is nil")
			return `{"error": "UI event dispatcher not set"}`
		}
		var payload struct {
			ID   string      `json:"id"`
			Name string      `json:"name"`
			Data interface{} `json:"data"`
		}
		if err := json.Unmarshal([]byte(jsonArgsPayload), &payload); err != nil {
			return fmt.Sprintf(`{"error": "Failed to parse UI event payload: %s"}`, err.Error())
		}
		fmt.Printf("Go: dispatching ui:event id=%s name=%s data=%v\n", payload.ID, payload.Name, payload.Data)
		if err := uiEventDispatcher(payload.ID, payload.Name, payload.Data); err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		return `{"status": "ok"}`
	}

	return `{"error": "Unsupported method call"}`
}

func (b *MobileBridge) PollNativeEvent() string {
	if globalAppInstance == nil {
		return ""
	}
	return globalAppInstance.Events.PollNativeEvent()
}

// HandleMessageFromFrontend is a package-level helper exposed to gomobile-generated Java wrappers.
func HandleMessageFromFrontend(methodKey string, jsonArgsPayload string) string {
	fmt.Printf("Go: HandleMessageFromFrontend method=%s payload=%s\n", methodKey, jsonArgsPayload)
	bridge := NewMobileBridge()
	return bridge.CallGoBackend(methodKey, jsonArgsPayload)
}

// HandleNativeAction is a package-level helper exposed to Java plugin packages.
func HandleNativeAction(methodKey string, jsonArgsPayload string) string {
	if globalAppInstance == nil {
		return `{"error": "Application core runtime context not active"}`
	}

	var rawArgs []json.RawMessage
	if err := json.Unmarshal([]byte(jsonArgsPayload), &rawArgs); err != nil {
		return fmt.Sprintf(`{"error": "Failed to extract arguments payload: %s"}`, err.Error())
	}

	result, err := globalAppInstance.InvokeNativeCall(methodKey, rawArgs)
	if err != nil {
		return fmt.Sprintf(`{"error": %q}`, err.Error())
	}

	responsePayload, _ := json.Marshal(map[string]interface{}{
		"result": result,
	})
	return string(responsePayload)
}
