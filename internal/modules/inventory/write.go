package inventory

import (
	"context"
	"fmt"
)

// Write dispatches write tool calls for the inventory module.
func (m *Module) Write(ctx context.Context, tool string, args map[string]any) (map[string]any, error) {
	switch tool {
	case "reset_all_stock":
		return m.resetAllStock(ctx)
	case "reset_category_stock":
		categoryID, _ := args["category_id"].(string)
		return m.resetCategoryStock(ctx, categoryID)
	case "reset_negative_stock":
		return m.resetNegativeStock(ctx)
	case "set_stock":
		variantID, _ := args["variant_id"].(string)
		storeID, _ := args["store_id"].(string)
		stock, _ := args["stock"].(float64)
		return m.setStock(ctx, variantID, storeID, stock)
	case "adjust_stock":
		itemID, _ := args["item_id"].(string)
		delta, _ := args["delta"].(float64)
		return m.adjustStock(ctx, itemID, delta)
	}
	return nil, fmt.Errorf("inventory: unknown write tool %q", tool)
}

func (m *Module) resetAllStock(ctx context.Context) (map[string]any, error) {
	ok, failed, err := m.client.ResetAllStock(ctx)
	if err != nil {
		return nil, fmt.Errorf("inventory: reset_all_stock: %w", err)
	}
	m.logger.Info("reset all stock", "updated", ok, "failed", failed)
	return map[string]any{"updated": ok, "failed": failed}, nil
}

func (m *Module) resetCategoryStock(ctx context.Context, categoryID string) (map[string]any, error) {
	if categoryID == "" {
		return nil, fmt.Errorf("inventory: reset_category_stock: category_id is required")
	}
	ok, failed, err := m.client.ResetCategoryStock(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("inventory: reset_category_stock: %w", err)
	}
	m.logger.Info("reset category stock", "category_id", categoryID, "updated", ok, "failed", failed)
	return map[string]any{"updated": ok, "failed": failed}, nil
}

func (m *Module) resetNegativeStock(ctx context.Context) (map[string]any, error) {
	ok, failed, err := m.client.ResetNegativeStock(ctx)
	if err != nil {
		return nil, fmt.Errorf("inventory: reset_negative_stock: %w", err)
	}
	m.logger.Info("reset negative stock", "updated", ok, "failed", failed)
	return map[string]any{"updated": ok, "failed": failed}, nil
}

func (m *Module) setStock(ctx context.Context, variantID, storeID string, stock float64) (map[string]any, error) {
	if variantID == "" {
		return nil, fmt.Errorf("inventory: set_stock: variant_id is required")
	}
	if storeID == "" {
		return nil, fmt.Errorf("inventory: set_stock: store_id is required")
	}
	if err := m.client.SetStock(ctx, variantID, storeID, stock); err != nil {
		return nil, fmt.Errorf("inventory: set_stock: %w", err)
	}
	return map[string]any{"variant_id": variantID, "store_id": storeID, "stock": stock}, nil
}

func (m *Module) adjustStock(ctx context.Context, itemID string, delta float64) (map[string]any, error) {
	if itemID == "" {
		return nil, fmt.Errorf("inventory: adjust_stock: item_id is required")
	}
	if err := m.client.AdjustStock(ctx, itemID, delta); err != nil {
		return nil, fmt.Errorf("inventory: adjust_stock: %w", err)
	}
	return map[string]any{"item_id": itemID, "delta": delta}, nil
}
