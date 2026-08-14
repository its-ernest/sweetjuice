package system

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

// SystemPlugin exposes OS-level information through a simple, stable API.
//
// It replaces the older osapi plugin with a clearer name and the same
// initial capability: querying basic system/device info from the native
// layer.
type SystemPlugin struct {
	app *core.Application
}

// SystemInfo represents basic OS/device information returned by the plugin.
type SystemInfo struct {
	// Common fields
	SystemName    string `json:"system_name"`
	SystemVersion string `json:"system_version"`
	Model         string `json:"model"`

	// Android specific
	SdkInt        int    `json:"sdk_int"`
	Release       string `json:"release"`
	Codename      string `json:"codename"`
	Manufacturer  string `json:"manufacturer"`
	Brand         string `json:"brand"`
	Board         string `json:"board"`
	Device        string `json:"device"`
	Product       string `json:"product"`
	Hardware      string `json:"hardware"`
	BaseOS        string `json:"base_os"`
	SecurityPatch string `json:"security_patch"`

	// iOS specific
	Name                string `json:"name"`
	LocalizedModel      string `json:"localized_model"`
	IdentifierForVendor string `json:"identifier_for_vendor"`
	IsPhysicalDevice    bool   `json:"is_physical_device"`
}

// NewPlugin creates a new instance of the System plugin.
func NewPlugin() *SystemPlugin {
	return &SystemPlugin{}
}

// Name returns the plugin name.
func (p *SystemPlugin) Name() string {
	return "system"
}

// Init initializes the plugin.
func (p *SystemPlugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

// GetInfo returns information about the current OS/device.
func (p *SystemPlugin) GetInfo() (SystemInfo, error) {
	var info SystemInfo
	result := core.CallNativePlatform("system:getInfo", "{}")

	if err := parseResultError(result); err != nil {
		return info, err
	}

	if err := json.Unmarshal([]byte(result), &info); err != nil {
		return info, fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	return info, nil
}

func parseResultError(result string) error {
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(result), &generic); err != nil {
		return nil
	}
	if errMsg, ok := generic["error"]; ok {
		return fmt.Errorf("native error: %v", errMsg)
	}
	return nil
}
