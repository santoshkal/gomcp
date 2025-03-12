// ./pkg/plugins/interfaces.go
package plugins

import (
	"context"
	"reflect"
)

// ToolHandler defines the signature that every plugin tool handler must implement.
type ToolHandler func(ctx context.Context, parameters map[string]interface{}) (interface{}, error)

// Tool is a common interface that every plugin tool should satisfy.
type Tool interface {
	// Name returns the tool's name.
	Name() string
	// Description returns a brief description.
	Description() string
	// Schema returns the JSON schema for input validation.
	Schema() map[string]interface{}
	// Handle executes the tool's functionality.
	Handle(ctx context.Context, parameters map[string]interface{}) (interface{}, error)
}

// Plugin defines the interface for a plugin that can register multiple tools.
type Plugin interface {
	// Tools returns a slice of tools provided by the plugin.
	Tools() []Tool
}

func HandlerSymbols() map[string]map[string]reflect.Value {
	return map[string]map[string]reflect.Value{
		"github.com/santoshkal/gomcp/pkg/plugins/plugins": {
			"ToolHandler": reflect.ValueOf((*ToolHandler)(nil)),
			"Tool":        reflect.ValueOf((*Tool)(nil)),
			"Plugin":      reflect.ValueOf((*Plugin)(nil)),
		},
	}
}
