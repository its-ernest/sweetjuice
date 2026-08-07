package calls

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

// CallsPlugin provides access to device call log.
type CallsPlugin struct {
	app *core.Application
}

// NewPlugin creates a new instance of the CallsPlugin.
func NewPlugin() *CallsPlugin {
	return &CallsPlugin{}
}

// Init registers the native callback for call log updates.
func (p *CallsPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("calls:changed", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(args[0], &payload); err != nil {
			return nil, err
		}

		app.Events.Emit("calls:changed", payload)
		return map[string]string{"status": "received"}, nil
	})

	return nil
}

// CallRecord represents a single call log entry.
type CallRecord struct {
	ID          int64  `json:"id"`
	Number      string `json:"number"`
	Type        string `json:"type"` // incoming, outgoing, missed, voicemail
	Date        int64  `json:"date"`
	Duration    int64  `json:"duration"`
	CachedName  string `json:"cached_name"`
	GeoLocation string `json:"geo_location,omitempty"`
}

// CallLog represents the device call log.
type CallLog struct {
	Calls    []CallRecord `json:"calls"`
	Count    int          `json:"count"`
}

// GetRecent fetches the most recent call log entries.
func (p *CallsPlugin) GetRecent(limit int) (CallLog, error) {
	var log CallLog
	payload, _ := json.Marshal(map[string]int{"limit": limit})
	result := core.CallNativePlatform("calls:getRecent", string(payload))

	if err := parsePluginError(result); err != nil {
		return log, err
	}

	if err := json.Unmarshal([]byte(result), &log); err != nil {
		return log, err
	}

	return log, nil
}

// GetAll fetches the entire call log.
func (p *CallsPlugin) GetAll() (CallLog, error) {
	var log CallLog
	result := core.CallNativePlatform("calls:getAll", "{}")

	if err := parsePluginError(result); err != nil {
		return log, err
	}

	if err := json.Unmarshal([]byte(result), &log); err != nil {
		return log, err
	}

	return log, nil
}

func (p *CallsPlugin) callNativePlatform(method string) (string, error) {
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
