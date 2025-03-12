// ./pkg/reg/servicestools.go
package reg

import (
	"bytes"
	"context"
	"fmt"
	"go/build"
	"os"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/unsafe"
	"gopkg.in/yaml.v2"

	"github.com/santoshkal/gomcp/pkg/mcp"
	"github.com/santoshkal/gomcp/pkg/plugins"
)

// Config represents the overall YAML configuration.
type Config struct {
	Services []ServiceConfig `yaml:"services"`
}

// ServiceConfig defines a service entry.
type ServiceConfig struct {
	Name    string       `yaml:"name"`
	Enabled bool         `yaml:"enabled"` // if false, skip this service
	Tools   []ToolConfig `yaml:"tools"`
}

// ToolConfig defines an individual tool.
type ToolConfig struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Enabled     bool                   `yaml:"enabled"` // if false, skip this tool
	Schema      map[string]interface{} `yaml:"schema"`
	Plugin      string                 `yaml:"plugin"` // Inline Go code for the handler.
}

// loadConfig reads and unmarshals the YAML file.
func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return &cfg, nil
}

// RegisterToolsFromConfig loads the configuration, evaluates each tool's script using Yaegi,
// and registers only the enabled tools using the provided Registry.
func RegisterToolsFromConfig(r mcp.Registry, configPath string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	for _, svc := range cfg.Services {
		if !svc.Enabled {
			continue
		}
		for _, tool := range svc.Tools {
			if !tool.Enabled {
				continue
			}
			goPath := build.Default.GOPATH
			fmt.Printf("GoPath: %v\n", goPath)

			// Create a new yaegi interpreter instance.
			var stdout, stderr bytes.Buffer
			i := interp.New(interp.Options{GoPath: goPath, Env: os.Environ(), Stdout: &stdout, Stderr: &stderr})
			if err := i.Use(stdlib.Symbols); err != nil {
				fmt.Printf("error loading package symbols: %v\n", err)
			}
			if err := i.Use(unsafe.Symbols); err != nil {
				fmt.Printf("error loading unsafe symbols: %v", err)
			}

			if err := i.Use(plugins.HandlerSymbols()); err != nil {
				fmt.Printf("error loading handler symbols: %v\n", err)
			}

			// code, err := os.ReadFile(tool.Script)
			// if err != nil {
			// 	fmt.Printf("error reading plugin file: %v", err)
			// }

			// Evaluate the provided script. The script must define a function "Handler".
			if _, err := i.Eval(fmt.Sprintf(`import "%s"`, tool.Plugin)); err != nil {
				return fmt.Errorf("failed to evaluate plugin for tool [%s]: %+v", tool.Name, err)
			}

			// if _, err := i.Eval(string(code)); err != nil {
			// 	fmt.Printf("error evaluating plugin code: %v", err)
			// }

			fmt.Printf("Script Path: %v\n ", tool.Plugin)

			// Retrieve the Handler symbol.
			v, err := i.Eval("Handler")
			if err != nil {
				return fmt.Errorf("failed to retrieve Handler symbol for tool %s: %v", tool.Name, err)
			}

			// Assert that the symbol has the correct signature.
			handler, ok := v.Interface().(func(context.Context, map[string]interface{}) (interface{}, error))
			if !ok {
				return fmt.Errorf("handler for tool %s does not have the correct signature", tool.Name)
			}

			// Register the tool.
			r.RegisterTool(tool.Name, tool.Description, tool.Schema, handler)
		}
	}
	return nil
}
