package actions

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	_titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	_selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	_normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("248"))

	_dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	_warnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))

	_errStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	_okStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42"))

	_keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	_hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// View renders the current actions state.
func (m Model) View() string {
	var content string
	switch m.state {
	case stateMenu:
		content = m.menuView()
	case statePicking:
		content = m.pickingView()
	case stateConfirm:
		content = m.confirmView()
	case stateRunning:
		content = m.runningView()
	case stateDone:
		content = m.doneView()
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) menuView() string {
	var sb strings.Builder
	sb.WriteString(_titleStyle.Render("ACCIONES") + "\n\n")
	for i, item := range _items {
		cursor := "  "
		style := _normalStyle
		if i == m.cursor {
			cursor = "▶ "
			style = _selectedStyle
		}
		sb.WriteString(cursor + style.Render(item.label) + "\n")
	}
	sb.WriteString("\n" + navHints(
		hint{"↑↓", "mover"},
		hint{"enter", "ejecutar"},
		hint{"esc", "volver"},
	))
	return sb.String()
}

func (m Model) pickingView() string {
	var sb strings.Builder
	sb.WriteString(_titleStyle.Render("ELEGIR CATEGORÍA") + "\n\n")

	if !m.catsLoaded {
		sb.WriteString(_hintStyle.Render("cargando categorías..."))
		return sb.String()
	}
	if len(m.cats) == 0 {
		sb.WriteString(_warnStyle.Render("no hay categorías disponibles"))
		sb.WriteString("\n\n" + navHints(hint{"esc", "volver"}))
		return sb.String()
	}

	for i, cat := range m.cats {
		cursor := "  "
		style := _normalStyle
		if i == m.cursor {
			cursor = "▶ "
			style = _selectedStyle
		}
		sb.WriteString(cursor + style.Render(cat.name) + "\n")
	}
	sb.WriteString("\n" + navHints(
		hint{"↑↓", "mover"},
		hint{"enter", "confirmar"},
		hint{"esc", "volver"},
	))
	return sb.String()
}

func (m Model) confirmView() string {
	var sb strings.Builder
	sb.WriteString(_titleStyle.Render("CONFIRMAR OPERACIÓN") + "\n\n")
	sb.WriteString(_warnStyle.Render(m.sel.label))
	if m.sel.needsCat {
		sb.WriteString("\n" + _normalStyle.Render("Categoría: "+m.selCat.name))
	}
	sb.WriteString("\n\n" + _dimStyle.Render("Esta operación no se puede deshacer."))
	sb.WriteString("\n\n" + navHints(
		hint{"y/enter", "confirmar"},
		hint{"n/esc", "cancelar"},
	))
	return sb.String()
}

func (m Model) runningView() string {
	return m.spinner.View() + "  " + _dimStyle.Render("Ejecutando "+m.sel.label+"...")
}

func (m Model) doneView() string {
	var sb strings.Builder
	if m.err != nil {
		sb.WriteString(_errStyle.Render("Error: " + m.err.Error()))
	} else {
		sb.WriteString(_okStyle.Render("✓ Operación completada") + "\n\n")
		sb.WriteString(fmt.Sprintf("  actualizados:  %s\n",
			_titleStyle.Render(fmt.Sprintf("%d", m.result.updated))))
		sb.WriteString(fmt.Sprintf("  fallidos:      %s",
			_warnStyle.Render(fmt.Sprintf("%d", m.result.failed))))
	}
	sb.WriteString("\n\n" + navHints(hint{"enter/esc", "volver al menú"}))
	return sb.String()
}

type hint struct{ key, label string }

func navHints(hints ...hint) string {
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = _keyStyle.Render("["+h.key+"]") + " " + _hintStyle.Render(h.label)
	}
	return strings.Join(parts, "  ")
}
