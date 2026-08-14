package notification_listener

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sweet-juice/sweetjuice/core"
)

// Notification represents a received or removed notification.
type Notification struct {
	PackageName string `json:"package_name"`
	ID          int    `json:"id"`
	Title       string `json:"title,omitempty"`
	Text        string `json:"text,omitempty"`
	IsOngoing   bool   `json:"is_ongoing"`
	Timestamp   int64  `json:"timestamp"`
}

// NotificationHandler is the signature for notification event callbacks.
type NotificationHandler func(Notification) error

// GrantedHandler is the signature for notification access granted callbacks.
type GrantedHandler func()

// NotificationListenerPlugin manages notification event handlers.
type NotificationListenerPlugin struct {
	app *core.Application
}

var globalRegistry = struct {
	mu      sync.RWMutex
	posted  map[string]NotificationHandler
	removed map[string]NotificationHandler
	granted map[string]GrantedHandler
}{
	posted:  make(map[string]NotificationHandler),
	removed: make(map[string]NotificationHandler),
	granted: make(map[string]GrantedHandler),
}

// NewPlugin creates a new NotificationListenerPlugin instance.
func NewPlugin() *NotificationListenerPlugin {
	return &NotificationListenerPlugin{}
}

// Init initializes the plugin and registers native method handlers.
func (p *NotificationListenerPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("notification-listener:posted", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var notif Notification
		if err := json.Unmarshal(args[0], &notif); err != nil {
			return nil, err
		}

		globalRegistry.mu.RLock()
		for _, handler := range globalRegistry.posted {
			if err := handler(notif); err != nil {
				globalRegistry.mu.RUnlock()
				return nil, err
			}
		}
		globalRegistry.mu.RUnlock()

		return map[string]string{"status": "processed"}, nil
	})

	app.RegisterNativeMethod("notification-listener:removed", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var notif Notification
		if err := json.Unmarshal(args[0], &notif); err != nil {
			return nil, err
		}

		globalRegistry.mu.RLock()
		defer globalRegistry.mu.RUnlock()

		for _, handler := range globalRegistry.removed {
			if err := handler(notif); err != nil {
				return nil, err
			}
		}

		return map[string]string{"status": "processed"}, nil
	})

	app.RegisterNativeMethod("notification-listener:granted", func(args []json.RawMessage) (interface{}, error) {
		globalRegistry.mu.RLock()
		for _, handler := range globalRegistry.granted {
			handler()
		}
		globalRegistry.mu.RUnlock()

		return map[string]string{"status": "processed"}, nil
	})

	return nil
}

// OnPosted registers a handler for posted notifications.
func (p *NotificationListenerPlugin) OnPosted(key string, handler NotificationHandler) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.posted[key] = handler
}

// OnRemoved registers a handler for removed notifications.
func (p *NotificationListenerPlugin) OnRemoved(key string, handler NotificationHandler) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.removed[key] = handler
}

// OnGranted registers a handler for when notification access is granted.
func (p *NotificationListenerPlugin) OnGranted(key string, handler GrantedHandler) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.granted[key] = handler
}
