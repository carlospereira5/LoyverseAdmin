package agent

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"
)

// Executor runs tool calls concurrently with panic recovery.
type Executor struct {
	registry *Registry
	logger   *log.Logger
}

// NewExecutor creates an Executor backed by reg.
func NewExecutor(reg *Registry, logger *log.Logger) *Executor {
	return &Executor{registry: reg, logger: logger}
}

// Execute runs all calls concurrently and returns one ToolResult per call.
// Individual tool errors are captured in the result — they never abort other calls.
func (e *Executor) Execute(ctx context.Context, calls []ToolCall) ([]ToolResult, error) {
	if len(calls) == 0 {
		return nil, nil
	}

	results := make([]ToolResult, len(calls))
	g, ctx := errgroup.WithContext(ctx)

	for i, call := range calls {
		g.Go(func() error {
			defer func() {
				if r := recover(); r != nil {
					e.logger.Error("panic in tool",
						"tool", call.Name,
						"panic", fmt.Sprintf("%v", r),
						"stack", string(debug.Stack()),
					)
					results[i] = ToolResult{
						Name:   call.Name,
						Result: map[string]any{"error": fmt.Sprintf("panic: %v", r)},
					}
				}
			}()

			start := time.Now()
			result, err := e.registry.Execute(ctx, call.Name, call.Args)
			elapsed := time.Since(start)
			if err != nil {
				e.logger.Error("tool failed", "tool", call.Name, "elapsed", elapsed, "err", err)
				results[i] = ToolResult{Name: call.Name, Result: map[string]any{"error": err.Error()}}
			} else {
				e.logger.Debug("tool completed", "tool", call.Name, "elapsed", elapsed)
				results[i] = ToolResult{Name: call.Name, Result: result}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("executor: %w", err)
	}
	return results, nil
}
