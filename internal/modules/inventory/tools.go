package inventory

import "LoyverseAdmin/internal/agent"

// ReadTools returns the read-only tool definitions for the inventory module.
func (m *Module) ReadTools() []agent.ToolDef {
	return []agent.ToolDef{
		{
			Name:        "list_inventory",
			Description: "Returns all stock levels across all variants and stores.",
		},
		{
			Name:        "get_product_stock",
			Description: "Returns the current stock for the first variant of a product.",
			Parameters: []agent.ParamDef{
				{Name: "item_id", Type: "string", Description: "Loyverse item ID (UUID)."},
			},
			Required: []string{"item_id"},
		},
	}
}

// WriteTools returns the write tool definitions for the inventory module.
func (m *Module) WriteTools() []agent.ToolDef {
	return []agent.ToolDef{
		{
			Name:        "reset_all_stock",
			Description: "Sets ALL stock levels to 0 across every variant and store. Irreversible.",
		},
		{
			Name:        "reset_category_stock",
			Description: "Sets all stock levels to 0 for every variant of every item in a specific category. Irreversible.",
			Parameters: []agent.ParamDef{
				{Name: "category_id", Type: "string", Description: "Loyverse category ID (UUID)."},
			},
			Required: []string{"category_id"},
		},
		{
			Name:        "reset_negative_stock",
			Description: "Sets all stock levels below zero back to 0.",
		},
		{
			Name:        "set_stock",
			Description: "Sets the absolute stock level for a specific variant and store.",
			Parameters: []agent.ParamDef{
				{Name: "variant_id", Type: "string", Description: "Loyverse variant ID (UUID)."},
				{Name: "store_id", Type: "string", Description: "Loyverse store ID (UUID)."},
				{Name: "stock", Type: "number", Description: "New absolute stock level."},
			},
			Required: []string{"variant_id", "store_id", "stock"},
		},
		{
			Name:        "adjust_stock",
			Description: "Adds a delta to the current stock of a product. Use negative delta to reduce.",
			Parameters: []agent.ParamDef{
				{Name: "item_id", Type: "string", Description: "Loyverse item ID (UUID)."},
				{Name: "delta", Type: "number", Description: "Units to add (positive) or subtract (negative)."},
			},
			Required: []string{"item_id", "delta"},
		},
	}
}
