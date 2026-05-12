package categories

import "LoyverseAdmin/internal/agent"

// ReadTools returns the read-only tool definitions for the categories module.
func (m *Module) ReadTools() []agent.ToolDef {
	return []agent.ToolDef{
		{
			Name:        "list_categories",
			Description: "Returns all product categories.",
		},
	}
}

// WriteTools returns the write tool definitions for the categories module.
func (m *Module) WriteTools() []agent.ToolDef {
	return []agent.ToolDef{
		{
			Name:        "create_category",
			Description: "Creates a new product category.",
			Parameters: []agent.ParamDef{
				{Name: "name", Type: "string", Description: "Category name."},
				{Name: "color", Type: "string", Description: "Category color (optional hex or name)."},
			},
			Required: []string{"name"},
		},
		{
			Name:        "update_category",
			Description: "Updates an existing category's name or color.",
			Parameters: []agent.ParamDef{
				{Name: "id", Type: "string", Description: "Category ID (UUID)."},
				{Name: "name", Type: "string", Description: "New category name."},
				{Name: "color", Type: "string", Description: "New category color (optional)."},
			},
			Required: []string{"id", "name"},
		},
		{
			Name:        "delete_category",
			Description: "Permanently deletes a category by ID.",
			Parameters: []agent.ParamDef{
				{Name: "id", Type: "string", Description: "Category ID (UUID)."},
			},
			Required: []string{"id"},
		},
	}
}
