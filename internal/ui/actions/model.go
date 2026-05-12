// Package actions implements the bulk-operations menu view.
package actions

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"LoyverseAdmin/internal/agent"
)

type state int

const (
	stateMenu state = iota
	statePicking
	stateConfirm
	stateRunning
	stateDone
)

type menuItem struct {
	tool     string
	label    string
	needsCat bool
}

var _items = []menuItem{
	{"reset_all_stock", "Resetear stock de todos los productos", false},
	{"reset_category_stock", "Resetear stock por categoría", true},
	{"reset_negative_stock", "Resetear stock negativo a 0", false},
	{"reset_all_costs", "Resetear costos de todos los productos", false},
	{"apply_standardized_names", "Estandarizar nombres (Title Case)", false},
}

// BackMsg is sent when the user wants to return to the dashboard.
type BackMsg struct{}

type category struct {
	id   string
	name string
}

type opResult struct {
	updated int
	failed  int
}

// Model is the Bubbletea model for the actions view.
type Model struct {
	state      state
	cursor     int
	cats       []category
	catsLoaded bool
	sel        menuItem
	selCat     category
	result     opResult
	err        error
	spinner    spinner.Model
	reg        *agent.Registry
	width      int
	height     int
}

// New creates an actions Model backed by the given registry.
func New(reg *agent.Registry) Model {
	return Model{spinner: newSpinner(), reg: reg}
}

// Reset returns the model cleared to stateMenu, preserving registry, categories, and dimensions.
func (m Model) Reset() Model {
	return Model{
		spinner:    newSpinner(),
		reg:        m.reg,
		cats:       m.cats,
		catsLoaded: m.catsLoaded,
		width:      m.width,
		height:     m.height,
	}
}

// Init fires the initial category fetch.
func (m Model) Init() tea.Cmd {
	return fetchCategories(m.reg)
}

// Update processes incoming messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case categoriesMsg:
		m.cats = msg.cats
		m.catsLoaded = true
		return m, nil

	case doneMsg:
		m.state = stateDone
		m.result = opResult{updated: msg.updated, failed: msg.failed}
		m.err = nil
		return m, nil

	case execErrMsg:
		m.state = stateDone
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.state == stateRunning {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.state {
	case stateMenu:
		return m.handleMenu(msg)
	case statePicking:
		return m.handlePicking(msg)
	case stateConfirm:
		return m.handleConfirm(msg)
	case stateDone:
		return m.handleDone(msg)
	}
	return m, nil
}

func (m Model) handleMenu(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(_items)-1 {
			m.cursor++
		}
	case "enter":
		m.sel = _items[m.cursor]
		if m.sel.needsCat {
			m.cursor = 0
			m.state = statePicking
		} else {
			m.state = stateConfirm
		}
	case "esc":
		return m, func() tea.Msg { return BackMsg{} }
	}
	return m, nil
}

func (m Model) handlePicking(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.cats)-1 {
			m.cursor++
		}
	case "enter":
		if len(m.cats) > 0 {
			m.selCat = m.cats[m.cursor]
			m.state = stateConfirm
		}
	case "esc":
		m.cursor = 0
		m.state = stateMenu
	}
	return m, nil
}

func (m Model) handleConfirm(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		m.state = stateRunning
		return m, tea.Batch(m.spinner.Tick, execute(m.reg, m.sel.tool, m.selCat.id))
	case "n", "esc":
		if m.sel.needsCat {
			m.state = statePicking
		} else {
			m.state = stateMenu
		}
	}
	return m, nil
}

func (m Model) handleDone(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter", " ":
		m.state = stateMenu
		m.cursor = 0
		m.err = nil
	}
	return m, nil
}

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	return s
}
