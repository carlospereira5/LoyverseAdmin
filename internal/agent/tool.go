package agent

// ToolDef defines a tool the agent can invoke.
type ToolDef struct {
	Name        string
	Description string
	Parameters  []ParamDef
	Required    []string
}

// ParamDef defines a single parameter of a tool.
type ParamDef struct {
	Name        string
	Type        string // "string", "integer", "number", "boolean"
	Description string
	Enum        []string
}

// ToolCall represents a tool invocation request.
type ToolCall struct {
	Name string
	Args map[string]any
}

// ToolResult holds the outcome of executing a tool.
type ToolResult struct {
	Name   string
	Result map[string]any
}
