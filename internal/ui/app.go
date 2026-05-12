// Package ui is the root Bubbletea application that routes between views.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"LoyverseAdmin/internal/agent"
	"LoyverseAdmin/internal/ui/actions"
	"LoyverseAdmin/internal/ui/dashboard"
)

type appView int

const (
	viewDash appView = iota
	viewActions
)

// App is the root Bubbletea model. It owns the active view and handles
// global key bindings (quit).
type App struct {
	view   appView
	dash   dashboard.Model
	act    actions.Model
	width  int
	height int
}

// New creates the root App model backed by the given tool registry.
func New(reg *agent.Registry) App {
	return App{
		dash: dashboard.New(reg),
		act:  actions.New(reg),
	}
}

// Init starts the active view.
func (a App) Init() tea.Cmd {
	return a.dash.Init()
}

// Update handles global keys and delegates everything else to the active view.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Window size — always forward to both models.
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		a.width, a.height = ws.Width, ws.Height
		a.dash, _ = a.dash.Update(ws)
		a.act, _ = a.act.Update(ws)
		return a, nil
	}

	// Global quit.
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return a, tea.Quit
	}

	// Actions model signals "go back".
	if _, ok := msg.(actions.BackMsg); ok {
		a.view = viewDash
		return a, nil
	}

	switch a.view {
	case viewDash:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "q":
				return a, tea.Quit
			case "a":
				a.view = viewActions
				a.act = a.act.Reset()
				return a, a.act.Init()
			}
		}
		var cmd tea.Cmd
		a.dash, cmd = a.dash.Update(msg)
		return a, cmd

	case viewActions:
		var cmd tea.Cmd
		a.act, cmd = a.act.Update(msg)
		return a, cmd
	}

	return a, nil
}

// View delegates rendering to the active view.
func (a App) View() string {
	switch a.view {
	case viewActions:
		return a.act.View()
	default:
		return a.dash.View()
	}
}
