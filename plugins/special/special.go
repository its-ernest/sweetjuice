package special

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

// SpecialType identifies the special permission being requested.
type SpecialType string

const (
	Accessibility    SpecialType = "accessibility"
	AllFilesAccess   SpecialType = "all_files_access"
	BatteryExemption SpecialType = "battery_exemption"
	NotificationAccess SpecialType = "notification_access"
)

// SpecialPlugin handles special Android permissions that require system-settings intents.
type SpecialPlugin struct {
	app *core.Application
}

// NewPlugin creates a new SpecialPlugin.
func NewPlugin() *SpecialPlugin {
	return &SpecialPlugin{}
}

// Init registers the native callback for permission status changes.
func (p *SpecialPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("special:permission:result", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var result map[string]interface{}
		if err := json.Unmarshal(args[0], &result); err != nil {
			return nil, err
		}

		app.Events.Emit("special:permission:changed", result)
		return map[string]string{"status": "processed"}, nil
	})

	return nil
}

// Request launches the system settings screen for the given special permission type.
// The user must manually enable the permission on that screen.
func (p *SpecialPlugin) Request(perm SpecialType) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"type": string(perm),
	})
	result := core.CallNativePlatform("special:request", string(payload))
	if err := parsePluginError(result); err != nil {
		return "", err
	}
	return result, nil
}

// Check returns whether the given special permission is currently granted.
func (p *SpecialPlugin) Check(perm SpecialType) (bool, error) {
	payload, _ := json.Marshal(map[string]string{
		"type": string(perm),
	})
	result := core.CallNativePlatform("special:check", string(payload))
	if err := parsePluginError(result); err != nil {
		return false, err
	}

	var response struct {
		Granted bool   `json:"granted"`
		Status  string `json:"status"`
	}
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		return false, err
	}
	return response.Granted, nil
}

// Convenience helpers.
func RequestAccessibility() (string, error)       { return NewPlugin().Request(Accessibility) }
func RequestAllFilesAccess() (string, error)      { return NewPlugin().Request(AllFilesAccess) }
func RequestBatteryExemption() (string, error)    { return NewPlugin().Request(BatteryExemption) }
func RequestNotificationAccess() (string, error)  { return NewPlugin().Request(NotificationAccess) }

func RequestAppSettings() (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"type": "app_settings",
	})
	result := core.CallNativePlatform("special:request", string(payload))
	if err := parsePluginError(result); err != nil {
		return "", err
	}
	return result, nil
}

func CheckAccessibility() (bool, error)      { return NewPlugin().Check(Accessibility) }
func CheckAllFilesAccess() (bool, error)     { return NewPlugin().Check(AllFilesAccess) }
func CheckBatteryExemption() (bool, error)   { return NewPlugin().Check(BatteryExemption) }
func CheckNotificationAccess() (bool, error) { return NewPlugin().Check(NotificationAccess) }

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
