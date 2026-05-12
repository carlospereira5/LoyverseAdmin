package main

import (
	"fmt"
	"os"

	"github.com/carlospereira5/loyverse"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"LoyverseAdmin/internal/agent"
	"LoyverseAdmin/internal/modules/categories"
	"LoyverseAdmin/internal/modules/inventory"
	"LoyverseAdmin/internal/modules/products"
	"LoyverseAdmin/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	logger := log.New(os.Stderr)

	token := os.Getenv("LOYVERSE_TOKEN")
	if token == "" {
		return fmt.Errorf("LOYVERSE_TOKEN environment variable is required")
	}

	client, err := loyverse.New(token)
	if err != nil {
		return fmt.Errorf("create loyverse client: %w", err)
	}

	reg := agent.NewRegistry()
	deps := agent.PortDeps{
		Client: client,
		Logger: logger,
	}

	modules := []agent.Module{
		products.New(),
		inventory.New(),
		categories.New(),
	}

	agent.Provision(reg, modules, deps)

	p := tea.NewProgram(ui.New(reg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI: %w", err)
	}
	return nil
}
