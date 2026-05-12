package inventory

import (
	"context"
	"fmt"

	"github.com/carlospereira5/loyverse"
)

// Read dispatches read tool calls for the inventory module.
func (m *Module) Read(ctx context.Context, tool string, args map[string]any) (map[string]any, error) {
	switch tool {
	case "list_inventory":
		return m.listInventory(ctx)
	case "get_product_stock":
		id, _ := args["item_id"].(string)
		return m.getProductStock(ctx, id)
	}
	return nil, fmt.Errorf("inventory: unknown read tool %q", tool)
}

func (m *Module) listInventory(ctx context.Context) (map[string]any, error) {
	levels, err := m.client.GetInventoryLevels(ctx)
	if err != nil {
		return nil, fmt.Errorf("inventory: list_inventory: %w", err)
	}
	return map[string]any{
		"levels": marshalLevels(levels),
		"total":  len(levels),
	}, nil
}

func (m *Module) getProductStock(ctx context.Context, itemID string) (map[string]any, error) {
	if itemID == "" {
		return nil, fmt.Errorf("inventory: get_product_stock: item_id is required")
	}
	stock, err := m.client.GetItemStock(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("inventory: get_product_stock: %w", err)
	}
	return map[string]any{"item_id": itemID, "in_stock": stock}, nil
}

func marshalLevels(levels []loyverse.InventoryLevel) []map[string]any {
	out := make([]map[string]any, len(levels))
	for i, l := range levels {
		out[i] = map[string]any{
			"variant_id": l.VariantID,
			"store_id":   l.StoreID,
			"in_stock":   l.InStock,
			"updated_at": l.UpdatedAt,
		}
	}
	return out
}
