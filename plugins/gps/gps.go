package gps

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

// GpsPlugin provides access to device GPS/location data.
type GpsPlugin struct {
	app *core.Application
}

// NewPlugin creates a new instance of the GpsPlugin.
func NewPlugin() *GpsPlugin {
	return &GpsPlugin{}
}

// Init registers the native callback for location updates.
func (p *GpsPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("gps:changed", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(args[0], &payload); err != nil {
			return nil, err
		}

		app.Events.Emit("gps:changed", payload)
		return map[string]string{"status": "received"}, nil
	})

	return nil
}

// Location represents a GPS location reading.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	Altitude  float64 `json:"altitude"`
	Speed     float64 `json:"speed"`
	Timestamp int64   `json:"timestamp"`
}

// GetCurrentLocation fetches the last known location from the device.
func (p *GpsPlugin) GetCurrentLocation() (Location, error) {
	var loc Location
	result := core.CallNativePlatform("gps:getCurrentLocation", "{}")

	if err := parsePluginError(result); err != nil {
		return loc, err
	}

	if err := json.Unmarshal([]byte(result), &loc); err != nil {
		return loc, err
	}

	return loc, nil
}

// StartMonitoring begins native GPS location monitoring.
func (p *GpsPlugin) StartMonitoring() (string, error) {
	return p.callNativePlatform("gps:startMonitoring")
}

// StopMonitoring disables native GPS location monitoring.
func (p *GpsPlugin) StopMonitoring() (string, error) {
	return p.callNativePlatform("gps:stopMonitoring")
}

func (p *GpsPlugin) callNativePlatform(method string) (string, error) {
	result := core.CallNativePlatform(method, "{}")
	if err := parsePluginError(result); err != nil {
		return result, err
	}
	return result, nil
}

func parsePluginError(result string) error {
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(result), &generic); err != nil {
		return nil
	}
	if errMsg, ok := generic["error"]; ok {
		return fmt.Errorf("%v", errMsg)
	}
	return nil
}
