package dashboard

import (
	"context"
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sync/errgroup"

	"LoyverseAdmin/internal/agent"
)

type loadedMsg struct{ stats Stats }
type errMsg struct{ err error }

// Fetch returns a tea.Cmd that loads all dashboard metrics concurrently.
func Fetch(reg *agent.Registry) tea.Cmd {
	return func() tea.Msg {
		var (
			mu     sync.Mutex
			result Stats
		)
		g, ctx := errgroup.WithContext(context.Background())

		g.Go(func() error {
			r, err := reg.Execute(ctx, "list_products", nil)
			if err != nil {
				return fmt.Errorf("list_products: %w", err)
			}
			var withoutImg int
			if items, ok := r["items"].([]map[string]any); ok {
				for _, item := range items {
					if url, _ := item["image_url"].(string); url == "" {
						withoutImg++
					}
				}
			}
			mu.Lock()
			result.TotalProducts = intVal(r, "total")
			result.WithoutImage = withoutImg
			mu.Unlock()
			return nil
		})

		g.Go(func() error {
			r, err := reg.Execute(ctx, "preview_standardized_names", nil)
			if err != nil {
				return fmt.Errorf("preview_standardized_names: %w", err)
			}
			mu.Lock()
			result.NonStandardNames = intVal(r, "total")
			mu.Unlock()
			return nil
		})

		g.Go(func() error {
			r, err := reg.Execute(ctx, "list_inventory", nil)
			if err != nil {
				return fmt.Errorf("list_inventory: %w", err)
			}
			var neg int
			if levels, ok := r["levels"].([]map[string]any); ok {
				for _, l := range levels {
					if stock, ok := l["in_stock"].(float64); ok && stock < 0 {
						neg++
					}
				}
			}
			mu.Lock()
			result.TotalLevels = intVal(r, "total")
			result.NegativeStock = neg
			mu.Unlock()
			return nil
		})

		g.Go(func() error {
			r, err := reg.Execute(ctx, "list_categories", nil)
			if err != nil {
				return fmt.Errorf("list_categories: %w", err)
			}
			mu.Lock()
			result.TotalCategories = intVal(r, "total")
			mu.Unlock()
			return nil
		})

		if err := g.Wait(); err != nil {
			return errMsg{err}
		}
		return loadedMsg{result}
	}
}

func intVal(m map[string]any, key string) int {
	v, _ := m[key].(int)
	return v
}
