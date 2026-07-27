// Package core provides the core runtime for Sweet Juice applications.
package core

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/sweet-juice/sweetjuice/ui"
)

// Options defines the application configuration.
type Options struct {
	Name    string
	Assets  embed.FS
	OnStart func(app *Application) error
}

type nativeMethod func([]json.RawMessage) (interface{}, error)

// Application represents the global lifecycle state container.
type Application struct {
	Name          string
	nativeMethods map[string]nativeMethod
	Events        *EventBus
	options       Options
}

// NewApplication creates a new instance of the Sweet Juice application.
func NewApplication(options Options) *Application {
	return &Application{
		Name:          options.Name,
		nativeMethods: make(map[string]nativeMethod),
		Events:        NewEventBus(),
		options:       options,
	}
}

// Run starts the application engine.
func (a *Application) Run() error {
	SetGlobalApp(a)

	// Register UI event dispatcher
	SetUIEventDispatcher(func(id string, event string, data interface{}) error {
		ui.DispatchEvent(id, event, data)
		return nil
	})

	if a.options.OnStart != nil {
		return a.options.OnStart(a)
	}

	return nil
}

// RegisterNativeMethod registers a Go function as a "Native Method".
func (a *Application) RegisterNativeMethod(methodKey string, fn func([]json.RawMessage) (interface{}, error)) {
	if a.nativeMethods == nil {
		a.nativeMethods = make(map[string]nativeMethod)
	}
	a.nativeMethods[methodKey] = fn
}

// InvokeNativeCall executes a previously registered native method.
func (a *Application) InvokeNativeCall(methodKey string, rawArgs []json.RawMessage) (interface{}, error) {
	if fn, exists := a.nativeMethods[methodKey]; exists {
		return fn(rawArgs)
	}
	return nil, fmt.Errorf("native method identity '%s' not registered with application", methodKey)
}
