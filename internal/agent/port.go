// Package agent defines the contracts that modules must implement and the
// dependencies provisioned to them at initialization time.
package agent

import (
	"context"

	"github.com/carlospereira5/loyverse"
	"github.com/charmbracelet/log"
)

// Module is the base interface every module must implement.
type Module interface {
	Name() string
	Init(deps PortDeps) error
}

// DataReader is a Module that exposes read-only tools.
type DataReader interface {
	Module
	ReadTools() []ToolDef
	Read(ctx context.Context, tool string, args map[string]any) (map[string]any, error)
}

// DataWriter is a Module that exposes write tools.
type DataWriter interface {
	Module
	WriteTools() []ToolDef
	Write(ctx context.Context, tool string, args map[string]any) (map[string]any, error)
}

// PortDeps holds the dependencies provisioned to each module during Init.
type PortDeps struct {
	Client *loyverse.Client
	Logger *log.Logger
}
