package plugins

import (
	"sync"
)

var (
	registry = make(map[string]Source)
	mu       sync.RWMutex
)

// Register adds a plugin to the global registry.
// This is typically called from init() in each plugin package.
func Register(s Source) {
	mu.Lock()
	defer mu.Unlock()
	registry[s.Name()] = s
}

// Get retrieves a plugin by name from the registry.
// Returns nil if the plugin is not registered.
func Get(name string) Source {
	mu.RLock()
	defer mu.RUnlock()
	return registry[name]
}

// List returns all registered plugin names.
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// Clear removes all plugins from the registry.
// Used primarily for testing.
func Clear() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[string]Source)
}
