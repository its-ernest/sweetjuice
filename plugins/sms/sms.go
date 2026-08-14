package sms

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

// SmsPlugin provides access to device SMS messages.
type SmsPlugin struct {
	app *core.Application
}

// NewPlugin creates a new instance of the SmsPlugin.
func NewPlugin() *SmsPlugin {
	return &SmsPlugin{}
}

// Init registers the native callback for SMS updates.
func (p *SmsPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("sms:changed", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(args[0], &payload); err != nil {
			return nil, err
		}

		app.Events.Emit("sms:changed", payload)
		return map[string]string{"status": "received"}, nil
	})

	return nil
}

// SmsMessage represents a single SMS message.
type SmsMessage struct {
	ID          int64  `json:"id"`
	Address     string `json:"address"`
	Body        string `json:"body"`
	Type        string `json:"type"` // inbox, sent, draft, etc.
	Timestamp   int64  `json:"timestamp"`
	Read        bool   `json:"read"`
}

// SmsFolder represents a collection of SMS messages.
type SmsFolder struct {
	Messages []SmsMessage `json:"messages"`
	Count    int          `json:"count"`
	Folder   string       `json:"folder"`
}

// GetRecent fetches the most recent SMS messages.
func (p *SmsPlugin) GetRecent(limit int) (SmsFolder, error) {
	var folder SmsFolder
	payload, _ := json.Marshal(map[string]int{"limit": limit})
	result := core.CallNativePlatform("sms:getRecent", string(payload))

	if err := parsePluginError(result); err != nil {
		return folder, err
	}

	if err := json.Unmarshal([]byte(result), &folder); err != nil {
		return folder, err
	}

	return folder, nil
}

// GetLast is an alias for GetRecent.
func (p *SmsPlugin) GetLast(limit int) (SmsFolder, error) {
	return p.GetRecent(limit)
}

// GetInbox fetches SMS messages from the inbox.
func (p *SmsPlugin) GetInbox() (SmsFolder, error) {
	var folder SmsFolder
	result := core.CallNativePlatform("sms:getInbox", "{}")

	if err := parsePluginError(result); err != nil {
		return folder, err
	}

	if err := json.Unmarshal([]byte(result), &folder); err != nil {
		return folder, err
	}

	return folder, nil
}

// GetSent fetches SMS messages from the sent folder.
func (p *SmsPlugin) GetSent() (SmsFolder, error) {
	var folder SmsFolder
	result := core.CallNativePlatform("sms:getSent", "{}")

	if err := parsePluginError(result); err != nil {
		return folder, err
	}

	if err := json.Unmarshal([]byte(result), &folder); err != nil {
		return folder, err
	}

	return folder, nil
}

// GetDrafts fetches SMS messages from the drafts folder.
func (p *SmsPlugin) GetDrafts() (SmsFolder, error) {
	var folder SmsFolder
	result := core.CallNativePlatform("sms:getDrafts", "{}")

	if err := parsePluginError(result); err != nil {
		return folder, err
	}

	if err := json.Unmarshal([]byte(result), &folder); err != nil {
		return folder, err
	}

	return folder, nil
}

// GetAll fetches all SMS messages.
func (p *SmsPlugin) GetAll() (SmsFolder, error) {
	var folder SmsFolder
	result := core.CallNativePlatform("sms:getAll", "{}")

	if err := parsePluginError(result); err != nil {
		return folder, err
	}

	if err := json.Unmarshal([]byte(result), &folder); err != nil {
		return folder, err
	}

	return folder, nil
}

func (p *SmsPlugin) callNativePlatform(method string) (string, error) {
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
