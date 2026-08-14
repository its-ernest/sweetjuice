// Package permission provides a standard plugin for handling Android runtime permissions.
// It bridges the Go application with the native Android permission system.
package permission

import (
	"encoding/json"
	"fmt"
	"sync"

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

type PermissionResult struct {
	Permission string `json:"permission"`
	Granted    bool   `json:"granted"`
}

type PermissionHandler func(PermissionResult)

var globalHandlers []PermissionHandler
var handlersMu sync.RWMutex

// Init initializes the plugin with the Sweet Juice application context and registers
// the "permissions:result" native callback handler.
func (p *PermissionPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("permissions:result", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var result PermissionResult
		if err := json.Unmarshal(args[0], &result); err != nil {
			fmt.Printf("permission: failed to unmarshal result: %v (raw: %s)\n", err, string(args[0]))
			return nil, err
		}

		fmt.Printf("permission: result received: %s granted=%v\n", result.Permission, result.Granted)
		app.Events.Emit("permissions:changed", result)

		handlersMu.RLock()
		for _, h := range globalHandlers {
			go h(result)
		}
		handlersMu.RUnlock()

		return map[string]string{"status": "processed"}, nil
	})

	return nil
}

// OnResult registers a callback for permission results.
func (p *PermissionPlugin) OnResult(h PermissionHandler) {
	handlersMu.Lock()
	defer handlersMu.Unlock()
	globalHandlers = append(globalHandlers, h)
}

func OnResult(h PermissionHandler) {
	NewPlugin().OnResult(h)
}

// Check queries the status of a specific permission synchronously.
func (p *PermissionPlugin) Check(permission string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"permission": permission,
	})
	result := core.CallNativePlatform("permissions:check", string(payload))
	if err := parsePluginError(result); err != nil {
		return "unknown", err
	}
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
	if err := parsePluginError(result); err != nil {
		return "", err
	}
	return result, nil
}

// RequestMultiple triggers native permission request dialogs on Android for multiple permissions.
// The native side requests them sequentially and reports individual results via permissions:changed.
func (p *PermissionPlugin) RequestMultiple(permissions []string) (string, error) {
	payload, _ := json.Marshal(map[string][]string{
		"permissions": permissions,
	})
	result := core.CallNativePlatform("permissions:requestMultiple", string(payload))
	if err := parsePluginError(result); err != nil {
		return "", err
	}
	return result, nil
}

func parsePluginError(result string) error {
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(result), &generic); err != nil {
		return nil
	}
	if errMsg, ok := generic["error"].(string); ok && errMsg != "" {
		return fmt.Errorf("%v", errMsg)
	}
	return nil
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
