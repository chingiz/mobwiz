package generator

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TemplateConfig represents the structure of templates.yaml
type TemplateConfig struct {
	Android map[string]ArchitectureConfig `yaml:"android"`
	Flutter map[string]ArchitectureConfig `yaml:"flutter"`
	IOS     map[string]ArchitectureConfig `yaml:"ios"`
}

// ArchitectureConfig represents the layers (domain, data, presentation)
type ArchitectureConfig map[string][]TemplateDefinition

// TemplateDefinition defines a single template file and its output path
type TemplateDefinition struct {
	OutputPath   string `yaml:"path"`     // Path template (e.g., "domain/models/{{pascalCase .Name}}.kt")
	TemplateFile string `yaml:"template"` // Template filename (e.g., "model.kt.tmpl")
}

// LoadTemplatesConfig loads the templates.yaml configuration
func LoadTemplatesConfig(embeddedFS fs.FS) (*TemplateConfig, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configPath := ""

	// 1. Check working directory and parents
	searchPaths := []string{wd}

	// 2. Check executable directory
	exe, err := os.Executable()
	if err == nil {
		searchPaths = append(searchPaths, filepath.Dir(exe))
	}

	for _, base := range searchPaths {
		curr := base
		for i := 0; i < 4; i++ {
			tryPath := filepath.Join(curr, "templates", "templates.yaml")
			if _, err := os.Stat(tryPath); err == nil {
				configPath = tryPath
				break
			}
			curr = filepath.Dir(curr)
			if curr == "." || curr == "/" {
				break
			}
		}
		if configPath != "" {
			break
		}
	}

	var content []byte
	if configPath != "" {
		content, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read templates.yaml: %w", err)
		}
	} else if embeddedFS != nil {
		content, err = fs.ReadFile(embeddedFS, "templates/templates.yaml")
		if err != nil {
			return nil, fmt.Errorf("templates.yaml not found in local or embedded filesystem")
		}
	} else {
		return nil, fmt.Errorf("templates.yaml not found")
	}

	var config TemplateConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse templates.yaml: %w", err)
	}

	return &config, nil
}

// LoadTemplate loads a template file content from disk
func LoadTemplate(platform, architecture, templateFile string, embeddedFS fs.FS) (string, error) {
	// Try to load from working directory first (development mode)
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if running from test_output subdirectory or similar
	// Look for templates in parent directories up to 3 levels

	// 1. Check working directory and parents
	searchPaths := []string{wd}

	// 2. Check executable directory
	exe, err := os.Executable()
	if err == nil {
		searchPaths = append(searchPaths, filepath.Dir(exe))
	}

	templatePath := ""
	for _, base := range searchPaths {
		curr := base
		for i := 0; i < 4; i++ {
			tryPath := filepath.Join(curr, "templates", platform, architecture, templateFile)
			if _, err := os.Stat(tryPath); err == nil {
				templatePath = tryPath
				break
			}
			curr = filepath.Dir(curr)
			if curr == "." || curr == "/" {
				break
			}
		}
		if templatePath != "" {
			break
		}
	}

	if templatePath != "" {
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return "", fmt.Errorf("failed to load template %s: %w", templatePath, err)
		}
		return string(content), nil
	}

	if embeddedFS != nil {
		embeddedPath := filepath.Join("templates", platform, architecture, templateFile)
		content, err := fs.ReadFile(embeddedFS, embeddedPath)
		if err == nil {
			return string(content), nil
		}
	}

	return "", fmt.Errorf("template not found: %s/%s/%s", platform, architecture, templateFile)
}

// GetTemplatesFor returns the list of templates for a specific platform and architecture,
// along with the normalized platform and architecture keys used to find them.
func (c *TemplateConfig) GetTemplatesFor(platform, architecture string) ([]TemplateDefinition, string, string, error) {
	platform = strings.ToLower(platform)
	// Handle platform aliases
	if strings.Contains(platform, "android") {
		platform = "android"
	} else if strings.Contains(platform, "flutter") {
		platform = "flutter"
	} else if strings.Contains(platform, "ios") {
		platform = "ios"
	}

	// Normalize architecture (simple mapping for now)
	architectureKey := ""
	if strings.Contains(strings.ToLower(architecture), "bloc") {
		architectureKey = "bloc"
	} else if strings.Contains(strings.ToLower(architecture), "mvvm") {
		architectureKey = "mvvm"
	} else {
		// Default fallback if simple matching fails, assume exact match or first available
		architectureKey = strings.ToLower(architecture)
	}

	var archConfig ArchitectureConfig
	var ok bool

	switch platform {
	case "android":
		archConfig, ok = c.Android[architectureKey]
	case "flutter":
		archConfig, ok = c.Flutter[architectureKey]
	case "ios":
		archConfig, ok = c.IOS[architectureKey]
	default:
		return nil, "", "", fmt.Errorf("unsupported platform: %s", platform)
	}

	if !ok {
		return nil, "", "", fmt.Errorf("architecture '%s' not found for platform '%s'", architecture, platform)
	}

	// Flatten the map (layers) into a single list of templates
	var templates []TemplateDefinition
	for _, layerTemplates := range archConfig {
		templates = append(templates, layerTemplates...)
	}

	return templates, platform, architectureKey, nil
}

// GetArchitectures returns the list of available architectures for a specific platform
func (c *TemplateConfig) GetArchitectures(platform string) []string {
	platform = strings.ToLower(platform)
	if strings.Contains(platform, "android") {
		platform = "android"
	} else if strings.Contains(platform, "flutter") {
		platform = "flutter"
	} else if strings.Contains(platform, "ios") {
		platform = "ios"
	}

	var archConfig map[string]ArchitectureConfig
	switch platform {
	case "android":
		archConfig = make(map[string]ArchitectureConfig)
		for k, v := range c.Android {
			archConfig[k] = v
		}
	case "flutter":
		archConfig = make(map[string]ArchitectureConfig)
		for k, v := range c.Flutter {
			archConfig[k] = v
		}
	case "ios":
		archConfig = make(map[string]ArchitectureConfig)
		for k, v := range c.IOS {
			archConfig[k] = v
		}
	default:
		return []string{}
	}

	var architectures []string
	for k := range archConfig {
		architectures = append(architectures, k)
	}
	return architectures
}
