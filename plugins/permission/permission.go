// Package permission provides a standard plugin for handling Android runtime permissions.
// It bridges the Go application with the native Android permission system.
package permission

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

// PermissionPlugin handles the Go-side logic for the permissions system.
// It allows Go code to check and request permissions via the native bridge.
type PermissionPlugin struct {
	app *core.Application
}

// NewPlugin creates a new instance of the PermissionPlugin.
func NewPlugin() *PermissionPlugin {
	return &PermissionPlugin{}
}

// Init initializes the plugin with the Sweet Juice application context and registers
// the "permissions:result" native callback handler.
func (p *PermissionPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("permissions:result", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var results []map[string]interface{}
		if err := json.Unmarshal(args[0], &results); err != nil {
			return nil, err
		}

		for _, result := range results {
			app.Events.Emit("permissions:changed", result)
		}
		return map[string]string{"status": "processed"}, nil
	})

	return nil
}

// Check queries the status of a specific permission synchronously.
func (p *PermissionPlugin) Check(permission string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"permission": permission,
	})
	result := core.CallNativePlatform("permissions:check", string(payload))
	var resultMap map[string]interface{}
	json.Unmarshal([]byte(result), &resultMap)
	if status, ok := resultMap["status"].(string); ok {
		return status, nil
	}
	return "unknown", nil
}

// Request triggers a native permission request dialog on Android.
func (p *PermissionPlugin) Request(permission string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"permission": permission,
	})
	result := core.CallNativePlatform("permissions:request", string(payload))
	return result, nil
}

// RequestMultiple triggers native permission request dialogs on Android for multiple permissions.
// The native side requests them sequentially and reports individual results via permissions:changed.
func (p *PermissionPlugin) RequestMultiple(permissions []string) (string, error) {
	payload, _ := json.Marshal(map[string][]string{
		"permissions": permissions,
	})
	result := core.CallNativePlatform("permissions:requestMultiple", string(payload))
	return result, nil
}

// Check queries the status of a specific permission synchronously.
func Check(permission string) (string, error) {
	return NewPlugin().Check(permission)
}

// Request triggers a native permission request dialog on Android.
func Request(permission string) (string, error) {
	return NewPlugin().Request(permission)
}

// RequestMultiple triggers native permission request dialogs on Android for multiple permissions.
func RequestMultiple(permissions []string) (string, error) {
	return NewPlugin().RequestMultiple(permissions)
}
