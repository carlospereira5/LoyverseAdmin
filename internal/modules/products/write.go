package products

import (
	"context"
	"fmt"
)

// Write dispatches write tool calls for the products module.
func (m *Module) Write(ctx context.Context, tool string, args map[string]any) (map[string]any, error) {
	switch tool {
	case "reset_all_costs":
		return m.resetAllCosts(ctx)
	case "apply_standardized_names":
		return m.applyStandardizedNames(ctx)
	case "update_product_name":
		id, _ := args["item_id"].(string)
		name, _ := args["name"].(string)
		return m.updateProductName(ctx, id, name)
	}
	return nil, fmt.Errorf("products: unknown write tool %q", tool)
}

func (m *Module) resetAllCosts(ctx context.Context) (map[string]any, error) {
	ok, failed, err := m.client.ResetAllCosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("products: reset_all_costs: %w", err)
	}
	m.logger.Info("reset all costs", "updated", ok, "failed", failed)
	return map[string]any{"updated": ok, "failed": failed}, nil
}

func (m *Module) applyStandardizedNames(ctx context.Context) (map[string]any, error) {
	items, err := m.client.GetItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("products: apply_standardized_names: get items: %w", err)
	}

	updates := make(map[string]string)
	skipped := 0
	for _, item := range items {
		proposed := standardizeName(item.Name)
		if proposed == item.Name {
			skipped++
			continue
		}
		updates[item.ID] = proposed
	}

	if len(updates) == 0 {
		return map[string]any{"updated": 0, "skipped": skipped, "failed": 0}, nil
	}

	ok, failed, err := m.client.UpdateItemNameBatch(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("products: apply_standardized_names: %w", err)
	}
	m.logger.Info("names standardized", "updated", ok, "skipped", skipped, "failed", failed)
	return map[string]any{"updated": ok, "skipped": skipped, "failed": failed}, nil
}

func (m *Module) updateProductName(ctx context.Context, id, name string) (map[string]any, error) {
	if id == "" {
		return nil, fmt.Errorf("products: update_product_name: item_id is required")
	}
	if name == "" {
		return nil, fmt.Errorf("products: update_product_name: name is required")
	}
	if err := m.client.UpdateItemName(ctx, id, name); err != nil {
		return nil, fmt.Errorf("products: update_product_name: %w", err)
	}
	return map[string]any{"item_id": id, "name": name}, nil
}
