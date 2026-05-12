package products

import (
	"context"
	"fmt"

	"github.com/carlospereira5/loyverse"
)

// Read dispatches read tool calls for the products module.
func (m *Module) Read(ctx context.Context, tool string, args map[string]any) (map[string]any, error) {
	switch tool {
	case "list_products":
		return m.listProducts(ctx)
	case "list_products_without_image":
		return m.listProductsWithoutImage(ctx)
	case "get_product":
		id, _ := args["item_id"].(string)
		return m.getProduct(ctx, id)
	case "preview_standardized_names":
		return m.previewStandardizedNames(ctx)
	}
	return nil, fmt.Errorf("products: unknown read tool %q", tool)
}

func (m *Module) listProducts(ctx context.Context) (map[string]any, error) {
	items, err := m.client.GetItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("products: list_products: %w", err)
	}
	return map[string]any{
		"items": marshalItems(items),
		"total": len(items),
	}, nil
}

func (m *Module) listProductsWithoutImage(ctx context.Context) (map[string]any, error) {
	items, err := m.client.GetItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("products: list_products_without_image: %w", err)
	}
	var missing []loyverse.Item
	for _, item := range items {
		if item.ImageURL == "" {
			missing = append(missing, item)
		}
	}
	return map[string]any{
		"items": marshalItems(missing),
		"total": len(missing),
	}, nil
}

func (m *Module) getProduct(ctx context.Context, id string) (map[string]any, error) {
	if id == "" {
		return nil, fmt.Errorf("products: get_product: item_id is required")
	}
	item, err := m.client.GetItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("products: get_product: %w", err)
	}
	return marshalItem(*item), nil
}

func (m *Module) previewStandardizedNames(ctx context.Context) (map[string]any, error) {
	items, err := m.client.GetItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("products: preview_standardized_names: %w", err)
	}
	type change struct {
		ID      string `json:"id"`
		Current string `json:"current"`
		Proposed string `json:"proposed"`
	}
	var changes []map[string]any
	for _, item := range items {
		proposed := standardizeName(item.Name)
		if proposed == item.Name {
			continue
		}
		changes = append(changes, map[string]any{
			"id":       item.ID,
			"current":  item.Name,
			"proposed": proposed,
		})
	}
	return map[string]any{
		"changes": changes,
		"total":   len(changes),
	}, nil
}

// ── marshal helpers ───────────────────────────────────────────────────────────

func marshalItem(item loyverse.Item) map[string]any {
	return map[string]any{
		"id":          item.ID,
		"name":        item.Name,
		"category_id": item.CategoryID,
		"track_stock": item.TrackStock,
		"image_url":   item.ImageURL,
		"cost":        item.Cost,
		"variants":    marshalVariants(item.Variants),
		"created_at":  item.CreatedAt,
		"updated_at":  item.UpdatedAt,
	}
}

func marshalItems(items []loyverse.Item) []map[string]any {
	out := make([]map[string]any, len(items))
	for i, item := range items {
		out[i] = marshalItem(item)
	}
	return out
}

func marshalVariants(variants []loyverse.Variant) []map[string]any {
	out := make([]map[string]any, len(variants))
	for i, v := range variants {
		out[i] = map[string]any{
			"id":            v.ID,
			"name":          v.Name,
			"sku":           v.SKU,
			"barcode":       v.Barcode,
			"default_price": v.DefaultPrice,
			"pricing_type":  v.PricingType,
			"cost":          v.Cost,
		}
	}
	return out
}
