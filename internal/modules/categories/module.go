// Package categories manages Loyverse product categories.
package categories

import (
	"github.com/carlospereira5/loyverse"
	"github.com/charmbracelet/log"

	"LoyverseAdmin/internal/agent"
)

var _ agent.DataReader = (*Module)(nil)
var _ agent.DataWriter = (*Module)(nil)

// Module provides read and write tools for Loyverse categories.
type Module struct {
	client *loyverse.Client
	logger *log.Logger
}

// New creates an uninitialized categories Module.
func New() *Module { return &Module{} }

// Name returns the module identifier.
func (m *Module) Name() string { return "categories" }

// Init wires the module's dependencies.
func (m *Module) Init(deps agent.PortDeps) error {
	m.client = deps.Client
	m.logger = deps.Logger
	return nil
}
