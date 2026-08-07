package core

import "sync"

// PluginDefinition describes a plugin to register on the native side.
type PluginDefinition struct {
	Name    string
	Domain  string
	JavaPkg string
	Class   string
}

var (
	pluginRegistryMu sync.RWMutex
	pluginRegistry   = make(map[string]*PluginDefinition)
)

// RegisterPlugin records a plugin definition so the native side can instantiate it.
func RegisterPlugin(def *PluginDefinition) {
	pluginRegistryMu.Lock()
	defer pluginRegistryMu.Unlock()
	pluginRegistry[def.Domain] = def
}

// GetRegisteredPlugins returns a snapshot of registered plugin definitions.
func GetRegisteredPlugins() []*PluginDefinition {
	pluginRegistryMu.RLock()
	defer pluginRegistryMu.RUnlock()
	out := make([]*PluginDefinition, 0, len(pluginRegistry))
	for _, p := range pluginRegistry {
		out = append(out, p)
	}
	return out
}
