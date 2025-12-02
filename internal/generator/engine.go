package generator

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/chingiz/mobwiz/internal/config"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Engine handles template rendering
type Engine struct {
	Config config.Config
}

// NewEngine creates a new template engine
func NewEngine(cfg config.Config) *Engine {
	return &Engine{Config: cfg}
}

// Render processes a template string with the given data
func (e *Engine) Render(tmplContent string, data interface{}) (string, error) {
	// Create a new template with Sprig functions
	// We also add a custom "pascalCase" helper if sprig doesn't have it exactly as we want,
	// but sprig has "CamelCase" which is usually PascalCase.
	// Let's alias some common ones to match the guide's {{snakeCase}} style.

	funcMap := sprig.TxtFuncMap()
	funcMap["snakeCase"] = sprig.TxtFuncMap()["snakecase"]
	funcMap["pascalCase"] = sprig.TxtFuncMap()["camelcase"] // Sprig's camelcase is actually PascalCase (e.g. "foo_bar" -> "FooBar")
	funcMap["camelCase"] = sprig.TxtFuncMap()["camelcase"]  // Just in case
	// Add custom function to convert package name to path (e.g. com.example -> com/example)
	funcMap["packagePath"] = func(packageName string) string {
		return filepath.ToSlash(filepath.Join(strings.Split(packageName, ".")...))
	}

	tmpl, err := template.New("module").Funcs(funcMap).Parse(tmplContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// RenderPath renders a file path which might contain template variables
func (e *Engine) RenderPath(pathTmpl string, data interface{}) (string, error) {
	return e.Render(pathTmpl, data)
}

// WriteFile writes content to a file, creating directories if needed
func (e *Engine) WriteFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return os.WriteFile(path, []byte(content), 0644)
}
