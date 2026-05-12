package categories

import (
	"context"
	"fmt"

	"github.com/carlospereira5/loyverse"
)

// Read dispatches read tool calls for the categories module.
func (m *Module) Read(ctx context.Context, tool string, args map[string]any) (map[string]any, error) {
	switch tool {
	case "list_categories":
		return m.listCategories(ctx)
	}
	return nil, fmt.Errorf("categories: unknown read tool %q", tool)
}

// Write dispatches write tool calls for the categories module.
func (m *Module) Write(ctx context.Context, tool string, args map[string]any) (map[string]any, error) {
	switch tool {
	case "create_category":
		name, _ := args["name"].(string)
		color, _ := args["color"].(string)
		return m.createCategory(ctx, name, color)
	case "update_category":
		id, _ := args["id"].(string)
		name, _ := args["name"].(string)
		color, _ := args["color"].(string)
		return m.updateCategory(ctx, id, name, color)
	case "delete_category":
		id, _ := args["id"].(string)
		return m.deleteCategory(ctx, id)
	}
	return nil, fmt.Errorf("categories: unknown write tool %q", tool)
}

func (m *Module) listCategories(ctx context.Context) (map[string]any, error) {
	cats, err := m.client.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("categories: list_categories: %w", err)
	}
	return map[string]any{
		"categories": marshalCategories(cats),
		"total":      len(cats),
	}, nil
}

func (m *Module) createCategory(ctx context.Context, name, color string) (map[string]any, error) {
	if name == "" {
		return nil, fmt.Errorf("categories: create_category: name is required")
	}
	cat, err := m.client.CreateOrUpdateCategory(ctx, loyverse.CategoryRequest{
		Name:  name,
		Color: color,
	})
	if err != nil {
		return nil, fmt.Errorf("categories: create_category: %w", err)
	}
	return marshalCategory(*cat), nil
}

func (m *Module) updateCategory(ctx context.Context, id, name, color string) (map[string]any, error) {
	if id == "" {
		return nil, fmt.Errorf("categories: update_category: id is required")
	}
	if name == "" {
		return nil, fmt.Errorf("categories: update_category: name is required")
	}
	cat, err := m.client.CreateOrUpdateCategory(ctx, loyverse.CategoryRequest{
		ID:    id,
		Name:  name,
		Color: color,
	})
	if err != nil {
		return nil, fmt.Errorf("categories: update_category: %w", err)
	}
	return marshalCategory(*cat), nil
}

func (m *Module) deleteCategory(ctx context.Context, id string) (map[string]any, error) {
	if id == "" {
		return nil, fmt.Errorf("categories: delete_category: id is required")
	}
	if err := m.client.DeleteCategory(ctx, id); err != nil {
		return nil, fmt.Errorf("categories: delete_category: %w", err)
	}
	return map[string]any{"deleted": id}, nil
}

// ── marshal helpers ───────────────────────────────────────────────────────────

func marshalCategory(cat loyverse.Category) map[string]any {
	return map[string]any{
		"id":         cat.ID,
		"name":       cat.Name,
		"color":      cat.Color,
		"created_at": cat.CreatedAt,
	}
}

func marshalCategories(cats []loyverse.Category) []map[string]any {
	out := make([]map[string]any, len(cats))
	for i, cat := range cats {
		out[i] = marshalCategory(cat)
	}
	return out
}
