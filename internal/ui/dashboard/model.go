// Package dashboard implements the audit dashboard view.
package dashboard

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"LoyverseAdmin/internal/agent"
)

// Stats holds the audit metrics shown on the dashboard.
type Stats struct {
	TotalProducts    int
	WithoutImage     int
	NonStandardNames int
	TotalLevels      int
	NegativeStock    int
	TotalCategories  int
}

// Model is the Bubbletea model for the dashboard view.
type Model struct {
	loading bool
	err     error
	stats   Stats
	spinner spinner.Model
	width   int
	height  int
	reg     *agent.Registry
}

// New creates a dashboard Model backed by the given registry.
func New(reg *agent.Registry) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	return Model{
		loading: true,
		spinner: s,
		reg:     reg,
	}
}

// Init starts the spinner animation and fires the initial data fetch.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, Fetch(m.reg))
}

// Update processes incoming messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case loadedMsg:
		m.loading = false
		m.err = nil
		m.stats = msg.stats
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "r" && !m.loading {
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, Fetch(m.reg))
		}
	}

	if m.loading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}
