package agent

import "context"

// Provision initializes all modules and registers their tools into reg.
// Modules that implement DataReader, DataWriter, or both are supported.
func Provision(reg *Registry, modules []Module, deps PortDeps) {
	for _, m := range modules {
		modDeps := deps
		modDeps.Logger = deps.Logger.WithPrefix(m.Name())

		if err := m.Init(modDeps); err != nil {
			deps.Logger.Error("module init failed", "module", m.Name(), "err", err)
			continue
		}

		var count int

		if r, ok := m.(DataReader); ok {
			for _, def := range r.ReadTools() {
				toolName := def.Name
				reader := r
				reg.Register(def, func(ctx context.Context, args map[string]any) (map[string]any, error) {
					return reader.Read(ctx, toolName, args)
				})
			}
			count += len(r.ReadTools())
		}

		if w, ok := m.(DataWriter); ok {
			for _, def := range w.WriteTools() {
				toolName := def.Name
				writer := w
				reg.Register(def, func(ctx context.Context, args map[string]any) (map[string]any, error) {
					return writer.Write(ctx, toolName, args)
				})
			}
			count += len(w.WriteTools())
		}

		modDeps.Logger.Info("provisioned", "tools", count)
	}
}
