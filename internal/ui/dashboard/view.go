package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	_titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	_dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	_cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(1, 2).
			MarginRight(2)

	_cardTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("248"))

	_numStyle = lipgloss.NewStyle().Bold(true)

	_warnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))

	_labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	_errStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	_keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	_hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// View renders the dashboard state.
func (m Model) View() string {
	switch {
	case m.loading:
		return loadingView(m)
	case m.err != nil:
		return errorView(m)
	default:
		return dashView(m)
	}
}

func loadingView(m Model) string {
	content := m.spinner.View() + "  Cargando datos de Loyverse..."
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func errorView(m Model) string {
	content := _errStyle.Render("Error al cargar: "+m.err.Error()) +
		"\n\n" + _hintStyle.Render("[r] reintentar  [q] salir")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func dashView(m Model) string {
	cw := cardWidth(m.width)
	sep := _dimStyle.Render(strings.Repeat("─", max(m.width-4, 10)))

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		catalogCard(m.stats, cw),
		inventoryCard(m.stats, cw),
		categoriesCard(m.stats, cw),
	)

	return strings.Join([]string{
		"",
		"  " + _titleStyle.Render("LoyverseAdmin"),
		"  " + sep,
		"",
		lipgloss.NewStyle().MarginLeft(2).Render(row),
		"",
		navHints(),
	}, "\n")
}

func catalogCard(s Stats, width int) string {
	var sb strings.Builder
	sb.WriteString(_cardTitleStyle.Render("CATÁLOGO") + "\n\n")
	sb.WriteString(statRow("productos", s.TotalProducts, false) + "\n")
	sb.WriteString(statRow("sin imagen", s.WithoutImage, s.WithoutImage > 0) + "\n")
	sb.WriteString(statRow("no estándar", s.NonStandardNames, s.NonStandardNames > 0))
	return _cardStyle.Width(width).Render(sb.String())
}

func inventoryCard(s Stats, width int) string {
	var sb strings.Builder
	sb.WriteString(_cardTitleStyle.Render("INVENTARIO") + "\n\n")
	sb.WriteString(statRow("niveles", s.TotalLevels, false) + "\n")
	sb.WriteString(statRow("negativos", s.NegativeStock, s.NegativeStock > 0))
	return _cardStyle.Width(width).Render(sb.String())
}

func categoriesCard(s Stats, width int) string {
	var sb strings.Builder
	sb.WriteString(_cardTitleStyle.Render("CATEGORÍAS") + "\n\n")
	sb.WriteString(statRow("categorías", s.TotalCategories, false))
	return _cardStyle.Width(width).Render(sb.String())
}

func statRow(label string, n int, warn bool) string {
	numStr := fmt.Sprintf("%5d", n)
	if warn && n > 0 {
		numStr = _warnStyle.Render(numStr)
	} else {
		numStr = _numStyle.Render(numStr)
	}
	return numStr + "  " + _labelStyle.Render(label)
}

func navHints() string {
	type hint struct{ key, label string }
	hints := []hint{
		{"a", "Acciones"},
		{"r", "Recargar"},
		{"q", "Salir"},
	}
	var parts []string
	for _, h := range hints {
		parts = append(parts, _keyStyle.Render("["+h.key+"]")+" "+_hintStyle.Render(h.label))
	}
	return "  " + strings.Join(parts, "  ")
}

func cardWidth(termWidth int) int {
	w := (termWidth - 12) / 3
	switch {
	case w < 20:
		return 20
	case w > 28:
		return 28
	default:
		return w
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
