package agent

import (
	"context"
	"fmt"
	"sync"
)

// Handler processes a tool call and returns a result map.
type Handler func(ctx context.Context, args map[string]any) (map[string]any, error)

// Registry stores tool definitions and their handlers.
// All methods are safe for concurrent use.
type Registry struct {
	mu       sync.RWMutex
	defs     []ToolDef
	handlers map[string]Handler
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]Handler)}
}

// Register adds a tool to the registry. An existing tool with the same name is
// replaced.
func (r *Registry) Register(def ToolDef, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, d := range r.defs {
		if d.Name == def.Name {
			r.defs[i] = def
			r.handlers[def.Name] = h
			return
		}
	}
	r.defs = append(r.defs, def)
	r.handlers[def.Name] = h
}

// Tools returns a copy of all registered ToolDefs.
func (r *Registry) Tools() []ToolDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ToolDef, len(r.defs))
	copy(out, r.defs)
	return out
}

// Execute dispatches a named tool call to its registered handler.
func (r *Registry) Execute(ctx context.Context, name string, args map[string]any) (map[string]any, error) {
	r.mu.RLock()
	h, ok := r.handlers[name]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown tool: %q", name)
	}
	return h(ctx, args)
}
