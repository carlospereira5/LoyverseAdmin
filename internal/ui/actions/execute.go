package actions

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"LoyverseAdmin/internal/agent"
)

type categoriesMsg struct{ cats []category }
type doneMsg struct{ updated, failed int }
type execErrMsg struct{ err error }

func fetchCategories(reg *agent.Registry) tea.Cmd {
	return func() tea.Msg {
		r, err := reg.Execute(context.Background(), "list_categories", nil)
		if err != nil {
			return execErrMsg{fmt.Errorf("cargar categorías: %w", err)}
		}
		raw, _ := r["categories"].([]map[string]any)
		cats := make([]category, 0, len(raw))
		for _, c := range raw {
			id, _ := c["id"].(string)
			name, _ := c["name"].(string)
			if id != "" {
				cats = append(cats, category{id: id, name: name})
			}
		}
		return categoriesMsg{cats}
	}
}

func execute(reg *agent.Registry, tool, categoryID string) tea.Cmd {
	return func() tea.Msg {
		var args map[string]any
		if categoryID != "" {
			args = map[string]any{"category_id": categoryID}
		}
		r, err := reg.Execute(context.Background(), tool, args)
		if err != nil {
			return execErrMsg{err}
		}
		updated, _ := r["updated"].(int)
		failed, _ := r["failed"].(int)
		return doneMsg{updated: updated, failed: failed}
	}
}
