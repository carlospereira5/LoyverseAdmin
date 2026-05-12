package products

import "LoyverseAdmin/internal/agent"

// ReadTools returns the read-only tool definitions for the products module.
func (m *Module) ReadTools() []agent.ToolDef {
	return []agent.ToolDef{
		{
			Name:        "list_products",
			Description: "Returns all products in the Loyverse catalog.",
		},
		{
			Name:        "list_products_without_image",
			Description: "Returns all products that have no image uploaded.",
		},
		{
			Name:        "get_product",
			Description: "Returns a single product by its Loyverse ID.",
			Parameters: []agent.ParamDef{
				{Name: "item_id", Type: "string", Description: "Loyverse item ID (UUID)."},
			},
			Required: []string{"item_id"},
		},
		{
			Name:        "preview_standardized_names",
			Description: "Shows which product names would change after applying Title Case standardization.",
		},
	}
}

// WriteTools returns the write tool definitions for the products module.
func (m *Module) WriteTools() []agent.ToolDef {
	return []agent.ToolDef{
		{
			Name:        "reset_all_costs",
			Description: "Sets the cost of every item and all its variants to 0. Irreversible.",
		},
		{
			Name:        "apply_standardized_names",
			Description: "Applies Title Case standardization to all product names that are not yet standardized.",
		},
		{
			Name:        "update_product_name",
			Description: "Updates the name of a single product.",
			Parameters: []agent.ParamDef{
				{Name: "item_id", Type: "string", Description: "Loyverse item ID (UUID)."},
				{Name: "name", Type: "string", Description: "New product name."},
			},
			Required: []string{"item_id", "name"},
		},
	}
}
