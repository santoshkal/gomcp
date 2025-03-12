// ./pkg/mcp/types.go
package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/santoshkal/gomcp/pkg/plugins"
)

const JSONRPCVersion = "2.0"

// RPCRequest defines the JSON-RPC request structure.
type RPCRequest struct {
	Version string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// RPCResponse defines the JSON-RPC response structure.
type RPCResponse struct {
	Version string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
	ID      *int            `json:"id"`
}

// RPCError defines an error in JSON-RPC responses.
type RPCError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// NewError creates a new RPCError with the given code and message.
func NewError(code int, msg string) *RPCError {
	return &RPCError{
		Code:    code,
		Message: msg,
	}
}

// String returns the string representation of the RPCError.
func (e *RPCError) String() string {
	return fmt.Sprintf("RPC Error [Code: %d]: %s", e.Code, e.Message)
}

// ToolCallArgs represents arguments for directly calling a tool.
type ToolCallArgs struct {
	ToolName   string                 `json:"tool_name"`
	Parameters map[string]interface{} `json:"parameters"`
}

// Registry defines the interface for registering tools.
type Registry interface {
	RegisterTool(name, description string, inputSchema map[string]interface{}, handler plugins.ToolHandler)
}
