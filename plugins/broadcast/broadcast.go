package broadcast

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sweet-juice/sweetjuice/core"
)

type BroadcastPlugin struct {
	app *core.Application
}

var (
	globalHandlers = make(map[string][]func(interface{}))
	handlersMu     sync.RWMutex
)

func NewPlugin() *BroadcastPlugin {
	return &BroadcastPlugin{}
}

func (p *BroadcastPlugin) Name() string {
	return "broadcast"
}

func (p *BroadcastPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("broadcast:post", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var payload struct {
			Name string      `json:"name"`
			Data interface{} `json:"data"`
		}

		if err := json.Unmarshal(args[0], &payload); err != nil {
			return nil, err
		}

		handlersMu.RLock()
		callbacks := globalHandlers[payload.Name]
		handlersMu.RUnlock()

		for _, cb := range callbacks {
			go cb(payload.Data)
		}

		return map[string]string{"status": "ok"}, nil
	})

	return nil
}

// On registers a handler for a specific Android Intent action.
// If it's a system broadcast (like BOOT_COMPLETED), ensure it's in your AndroidManifest.xml.
func (p *BroadcastPlugin) On(action string, callback func(interface{})) {
	handlersMu.Lock()
	globalHandlers[action] = append(globalHandlers[action], callback)
	handlersMu.Unlock()

	// Notify native side to register a dynamic receiver for this action
	payload, _ := json.Marshal(map[string]string{
		"action": action,
	})
	core.CallNativePlatform("broadcast:register", string(payload))
}

func On(action string, callback func(interface{})) {
	NewPlugin().On(action, callback)
}

// Post sends a system-wide broadcast (Intent) on the native side.
func (p *BroadcastPlugin) Post(action string, extras map[string]interface{}) {
	payload, _ := json.Marshal(map[string]interface{}{
		"action": action,
		"extras": extras,
	})
	core.CallNativePlatform("broadcast:send", string(payload))
}

func Post(action string, extras map[string]interface{}) {
	NewPlugin().Post(action, extras)
}
